# SAMS MCP Server

This directory contains a Model Context Protocol (MCP) server for the Smart Asset Management System (SAMS).
This server exposes tools that an AI agent (like Claude) can use to query the SAMS backend API and retrieve information about assets.

## Prerequisites

- Python 3.10+
- `uv` package manager (https://github.com/astral-sh/uv)

## Setup

1.  **Navigate to this directory**:
    ```bash
    cd mcp-server
    ```

2.  **Install dependencies**:
    If you haven't already, install the required Python packages using `uv`.
    ```bash
    uv add "mcp[cli]" httpx
    ```

## Running the Server

The MCP server communicates over `stdio`. You can run it directly using `uv`:

```bash
uv run python server.py
```

The server will start and wait for JSON-RPC messages from an MCP client.

## Available Tools

-   **`get_asset_summary()`**:
    -   Description: Retrieves a summary of all assets, including total count, total value, and counts of active/critical assets.
    -   Usage: `get_asset_summary`

-   **`get_recent_assets(limit: int = 5)`**:
    -   Description: Retrieves a list of the most recently added assets.
    -   Usage: `get_recent_assets` or `get_recent_assets --limit 10`

## Connecting to a Client (e.g., Claude for Desktop)

To connect this server to an MCP client like Claude for Desktop, you need to edit the client's configuration file.

1.  **Find the configuration file**: For Claude for Desktop on macOS, it's located at `~/Library/Application Support/Claude/claude_desktop_config.json`.

2.  **Add the server configuration**: Add the following JSON object to the `mcpServers` key. Make sure to replace `/ABSOLUTE/PATH/TO/SAMS/mcp-server` with the actual absolute path to this directory on your machine.

    ```json
    {
      "mcpServers": {
        "sams_info_provider": {
          "command": "/ABSOLUTE/PATH/TO/SAMS/mcp-server/.venv/bin/python",
          "args": [
            "/ABSOLUTE/PATH/TO/SAMS/mcp-server/server.py"
          ],
          "workingDirectory": "/ABSOLUTE/PATH/TO/SAMS/mcp-server"
        }
      }
    }
    ```
    
    *You can find the absolute path by navigating to this directory in your terminal and running `pwd`.*

3.  **Restart the client**: After saving the configuration, completely restart the client application for the changes to take effect.
