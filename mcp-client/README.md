# SAMS MCP Client

This directory contains a client application that connects to the SAMS MCP Server and uses the Google Gemini API to provide a natural language interface for querying asset information.

## Prerequisites

- Python 3.10+
- `uv` package manager (https://github.com/astral-sh/uv)
- A Google Gemini API Key.

## Setup

1.  **Navigate to this directory**:
    ```bash
    cd mcp-client
    ```

2.  **Install dependencies**:
    If you haven't already, install the required Python packages using `uv`.
    ```bash
    uv add "google-genai" "mcp[cli]" "python-dotenv"
    ```

3.  **Set up your environment variables**:
    Copy the `.env.example` file to a new file named `.env`.
    ```bash
    cp .env.example .env
    ```
    Open the `.env` file and add your Google Gemini API key.
    ```
    GEMINI_API_KEY="YOUR_API_KEY_HERE"
    ```

## Running the Client

First, make sure your main SAMS backend is running, as the MCP server will call its API endpoints.

To run the client, simply execute the `client.py` script using `uv`:

```bash
uv run python client.py
```

The script will:
1.  Automatically start the `mcp-server` in the background.
2.  Connect to it and retrieve the available tools.
3.  Provide a chat prompt where you can ask questions about your assets.

### Example Questions

-   `give me a summary of my assets`
-   `show me the 3 most recent assets`
-   `what are the latest assets?`

Type `exit` to quit the application.
