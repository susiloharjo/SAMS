from typing import Any
import httpx
from mcp.server.fastmcp import FastMCP
import logging
from fastapi import FastAPI, HTTPException
import uvicorn
import inspect

# Configure logging to stderr
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

# --- MCP Tool Definitions ---

mcp = FastMCP(
    name="sams_info_provider",
    # title="SAMS Information Provider", # Unexpected keyword argument
    # description="Provides tools to query the Smart Asset Management System (SAMS)."
)

SAMS_API_BASE = "http://sams-backend:8080/api/v1" # Use Docker service name
USER_AGENT = "sams-mcp-server/1.0"

async def make_sams_request(endpoint: str) -> dict[str, Any] | None:
    """Make a request to the SAMS API with proper error handling."""
    url = f"{SAMS_API_BASE}{endpoint}"
    headers = {"User-Agent": USER_AGENT}
    async with httpx.AsyncClient() as client:
        try:
            logging.info(f"Making request to SAMS API: {url}")
            response = await client.get(url, headers=headers, timeout=10.0)
            response.raise_for_status()
            logging.info(f"Request to {url} successful with status {response.status_code}")
            return response.json()
        except httpx.RequestError as exc:
            logging.error(f"An error occurred while requesting {exc.request.url!r}: {exc!r}")
            return None
        except httpx.HTTPStatusError as exc:
            logging.error(f"Error response {exc.response.status_code} while requesting {exc.request.url!r}: {exc.response.text!r}")
            return None

@mcp.tool()
async def get_asset_summary() -> str:
    """
    Retrieves a summary of all assets in the SAMS database.
    This includes total number of assets, total value, number of active assets, and number of critical assets.
    """
    logging.info("Executing get_asset_summary tool")
    data = await make_sams_request("/assets/summary")

    if not data or not data.get("data"):
        return "Unable to fetch asset summary from SAMS API."

    summary = data["data"]
    return (
        f"Asset Summary:\n"
        f"- Total Assets: {summary.get('total_assets', 'N/A')}\n"
        f"- Total Value: ${summary.get('total_value', 0):,.2f}\n"
        f"- Active Assets: {summary.get('active_assets', 'N/A')}\n"
        f"- Critical Assets: {summary.get('critical_assets', 'N/A')}"
    )

@mcp.tool()
async def get_recent_assets(limit: int = 5) -> str:
    """
    Retrieves a list of the most recently added assets.

    Args:
        limit: The number of recent assets to retrieve. Defaults to 5.
    """
    logging.info(f"Executing get_recent_assets tool with limit={limit}")
    data = await make_sams_request(f"/assets?limit={limit}&page=1")

    if not data or "data" not in data:
        return "Unable to fetch recent assets from SAMS API."

    assets = data["data"]
    if not assets:
        return "No recent assets found."

    asset_list = []
    for asset in assets:
        asset_list.append(
            f"- ID: {asset.get('id')}\n"
            f"  Name: {asset.get('name', 'N/A')}\n"
            f"  Status: {asset.get('status', 'N/A')}\n"
            f"  Acquisition Date: {asset.get('acquisition_date', 'N/A')}"
        )
    
    return "Most Recent Assets:\n" + "\n".join(asset_list)


# --- FastAPI HTTP Server ---

app = FastAPI()

@app.get("/tools")
async def list_tools():
    """Endpoint to get the list of available tools and their schemas."""
    # This part is conceptual and depends on how FastMCP stores tool schemas.
    # We introspect the mcp object to build the schema.
    tool_schemas = []
    for name, tool_func in mcp.tools.items():
        sig = inspect.signature(tool_func)
        schema = {
            "name": name,
            "description": tool_func.__doc__.strip() if tool_func.__doc__ else "No description.",
            "parameters": {
                "type": "object",
                "properties": {},
                "required": [],
            }
        }
        for param in sig.parameters.values():
            if param.name == 'self': continue # Skip self
            param_type = "string" # Default, can be improved with type hints
            if param.annotation == int: param_type = "integer"
            if param.annotation == float: param_type = "number"
            if param.annotation == bool: param_type = "boolean"
            
            schema["parameters"]["properties"][param.name] = {"type": param_type, "description": ""}
            if param.default is inspect.Parameter.empty:
                schema["parameters"]["required"].append(param.name)
        tool_schemas.append(schema)

    return {"tools": tool_schemas}

@app.post("/call/{tool_name}")
async def call_tool(tool_name: str, params: dict[str, Any]):
    """Endpoint to execute a specific tool with given parameters."""
    if tool_name not in mcp.tools:
        raise HTTPException(status_code=404, detail=f"Tool '{tool_name}' not found.")

    tool_func = mcp.tools[tool_name]
    try:
        logging.info(f"Calling tool '{tool_name}' with params: {params}")
        # Note: This assumes all tool functions are async
        result = await tool_func(**params)
        return {"result": result}
    except Exception as e:
        logging.error(f"Error executing tool '{tool_name}': {e}", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))


if __name__ == "__main__":
    logging.info("Starting SAMS MCP server as an HTTP server on port 8081")
    uvicorn.run(app, host="0.0.0.0", port=8081)
