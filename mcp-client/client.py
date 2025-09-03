import asyncio
import json
import os
import subprocess
import logging
import google.generativeai as genai
from google.generativeai import types
from dotenv import load_dotenv

# --- Configuration ---
# Load the .env file from the project root directory
dotenv_path = os.path.join(os.path.dirname(__file__), '..', '.env')
load_dotenv(dotenv_path=dotenv_path)

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
GEMINI_API_KEY = os.environ.get("GEMINI_API_KEY")
MCP_SERVER_PATH = os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'mcp-server'))
MCP_VENV_PYTHON = os.path.join(MCP_SERVER_PATH, '.venv', 'bin', 'python')
MCP_SERVER_SCRIPT = os.path.join(MCP_SERVER_PATH, 'server.py')


class MCPClient:
    """A conceptual client to communicate with an MCP server over stdio."""
    def __init__(self, command, args):
        self.command = command
        self.args = args
        self.process = None
        self.reader = None
        self.writer = None
        self.request_id = 0
        self.tool_schemas = []

    async def start(self):
        """Starts the MCP server subprocess."""
        logging.info(f"Starting MCP server with: {' '.join([self.command] + self.args)}")
        self.process = await asyncio.create_subprocess_exec(
            self.command,
            *self.args,
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        self.reader = self.process.stdout
        self.writer = self.process.stdin
        # Start a task to log stderr from the server
        asyncio.create_task(self._log_stderr())
        await self._initialize()

    async def _log_stderr(self):
        """Logs stderr from the subprocess for debugging."""
        while self.process.stderr.at_eof() is False:
            line = await self.process.stderr.readline()
            if line:
                logging.warning(f"[MCP Server stderr] {line.decode().strip()}")

    async def _send_request(self, method, params=None):
        """Sends a JSON-RPC request to the server."""
        self.request_id += 1
        request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": method,
            "params": params or {}
        }
        message = json.dumps(request)
        logging.info(f"Sending to MCP server: {message}")
        self.writer.write(message.encode('utf-8') + b'\n')
        await self.writer.drain()

    async def _read_response(self):
        """Reads a JSON-RPC response from the server."""
        # Note: A real implementation would need to handle buffered reading
        # and splitting of multiple JSON objects. This is a simplification.
        line = await self.reader.readline()
        if not line:
            return None
        response = json.loads(line.decode('utf-8'))
        logging.info(f"Received from MCP server: {response}")
        return response

    async def _initialize(self):
        """Initializes the connection and gets tool schemas."""
        await self._send_request("initialize", {"protocolVersion": "1.0"})
        response = await self._read_response()
        if response and "result" in response and "tools" in response["result"]:
            self.tool_schemas = response["result"]["tools"]
            logging.info(f"Successfully initialized and received {len(self.tool_schemas)} tools.")
        else:
            logging.error("Failed to initialize or get tools from MCP server.")
            raise ConnectionError("Could not initialize MCP server.")

    async def execute_tool(self, tool_name: str, parameters: dict):
        """Executes a tool on the MCP server."""
        await self._send_request("executeTool", {"name": tool_name, "parameters": parameters})
        response = await self._read_response()
        if response and "result" in response:
            return response["result"]
        return {"error": "Failed to execute tool or get a valid response."}

    async def stop(self):
        """Stops the MCP server subprocess."""
        if self.process:
            self.process.terminate()
            await self.process.wait()
            logging.info("MCP server stopped.")


async def run_chat_loop():
    """Main chat loop to interact with Gemini and the MCP server."""
    if not GEMINI_API_KEY:
        print("GEMINI_API_KEY not found in environment variables. Please set it in a .env file.")
        return

    genai.configure(api_key=GEMINI_API_KEY)
    mcp_client = MCPClient(MCP_VENV_PYTHON, [MCP_SERVER_SCRIPT])
    
    try:
        await mcp_client.start()

        # Convert MCP tool schemas to Gemini's format
        gemini_tools = []
        for schema in mcp_client.tool_schemas:
            gemini_tools.append(types.FunctionDeclaration(
                name=schema['name'],
                description=schema['description'],
                parameters=types.Schema(
                    type=types.Type.OBJECT,
                    properties={
                        p['name']: types.Schema(type=types.Type.STRING) # Simplified: assuming all params are strings
                        for p in schema.get('parameters', [])
                    },
                    required=[p['name'] for p in schema.get('parameters', []) if p.get('required')]
                )
            ))
        
        model = genai.GenerativeModel(
            model_name='gemini-1.5-flash',
            tools=gemini_tools
        )
        chat = model.start_chat(enable_automatic_function_calling=False)

        print("\n--- SAMS AI Assistant ---")
        print("Ask me about your assets! (Type 'exit' to quit)")
        
        while True:
            user_input = input("> ")
            if user_input.lower() == 'exit':
                break

            response = await chat.send_message_async(user_input)
            
            while response.candidates[0].function_calls:
                function_calls = response.candidates[0].function_calls
                tool_results = []
                
                for call in function_calls:
                    tool_name = call.name
                    params = dict(call.args)
                    print(f"🤖 Gemini wants to call tool: {tool_name}({params})")
                    
                    # Execute the tool via MCP client
                    mcp_result = await mcp_client.execute_tool(tool_name, params)

                    tool_results.append(types.Tool.FunctionResponse(
                        name=tool_name,
                        response=mcp_result
                    ))

                # Send the tool results back to Gemini
                response = await chat.send_message_async(
                    types.Content(parts=[types.Part.from_tool_response(tool_response=tool_results)])
                )

            print(f"Assistant: {response.text}")

    except ConnectionError as e:
        logging.error(e)
    except Exception as e:
        logging.error(f"An unexpected error occurred: {e}", exc_info=True)
    finally:
        await mcp_client.stop()


if __name__ == "__main__":
    asyncio.run(run_chat_loop())
