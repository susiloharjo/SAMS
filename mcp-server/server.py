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

SAMS_API_BASE = "http://sams-backend:8080/api/v1" # Use Docker service name for container communication
USER_AGENT = "sams-mcp-server/1.0"

async def make_sams_request(endpoint: str, method: str = "GET", data: dict = None) -> dict[str, Any] | None:
    """Make a request to the SAMS API with proper error handling."""
    url = f"{SAMS_API_BASE}{endpoint}"
    headers = {"User-Agent": USER_AGENT, "Content-Type": "application/json"}
    
    async with httpx.AsyncClient() as client:
        try:
            logging.info(f"Making {method} request to SAMS API: {url}")
            if method == "GET":
                response = await client.get(url, headers=headers, timeout=10.0)
            elif method == "POST":
                response = await client.post(url, headers=headers, json=data, timeout=10.0)
            else:
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

@mcp.tool()
async def search_assets(query: str, limit: int = 10) -> str:
    """
    Search for assets by name, serial number, model, or description.

    Args:
        query: The search term to look for in asset names, serial numbers, models, or descriptions.
        limit: Maximum number of results to return. Defaults to 10.
    """
    logging.info(f"Executing search_assets tool with query='{query}' and limit={limit}")
    data = await make_sams_request(f"/assets?search={query}&limit={limit}&page=1")

    if not data or "data" not in data:
        return f"Unable to search assets with query '{query}' from SAMS API."

    assets = data["data"]
    if not assets:
        return f"No assets found matching the query '{query}'."

    asset_list = []
    for asset in assets:
        asset_list.append(
            f"- ID: {asset.get('id')}\n"
            f"  Name: {asset.get('name', 'N/A')}\n"
            f"  Serial Number: {asset.get('serial_number', 'N/A')}\n"
            f"  Model: {asset.get('model', 'N/A')}\n"
            f"  Status: {asset.get('status', 'N/A')}\n"
            f"  Category: {asset.get('category', {}).get('name', 'N/A')}\n"
            f"  Department: {asset.get('department', {}).get('name', 'N/A')}"
        )
    
    return f"Search Results for '{query}' ({len(assets)} assets found):\n" + "\n".join(asset_list)

@mcp.tool()
async def get_assets_by_status(status: str, limit: int = 20) -> str:
    """
    Get assets filtered by their status (Active, Inactive, Maintenance, Disposed).

    Args:
        status: The status to filter by (Active, Inactive, Maintenance, Disposed).
        limit: Maximum number of results to return. Defaults to 20.
    """
    logging.info(f"Executing get_assets_by_status tool with status='{status}' and limit={limit}")
    data = await make_sams_request(f"/assets?status={status}&limit={limit}&page=1")

    if not data or "data" not in data:
        return f"Unable to fetch assets with status '{status}' from SAMS API."

    assets = data["data"]
    if not assets:
        return f"No assets found with status '{status}'."

    asset_list = []
    for asset in assets:
        asset_list.append(
            f"- ID: {asset.get('id')}\n"
            f"  Name: {asset.get('name', 'N/A')}\n"
            f"  Serial Number: {asset.get('serial_number', 'N/A')}\n"
            f"  Category: {asset.get('category', {}).get('name', 'N/A')}\n"
            f"  Department: {asset.get('department', {}).get('name', 'N/A')}\n"
            f"  Current Value: ${asset.get('current_value', 0):,.2f}"
        )
    
    return f"Assets with Status '{status}' ({len(assets)} assets found):\n" + "\n".join(asset_list)

@mcp.tool()
async def get_assets_by_category(category: str, limit: int = 20) -> str:
    """
    Get assets filtered by their category.

    Args:
        category: The category name to filter by.
        limit: Maximum number of results to return. Defaults to 20.
    """
    logging.info(f"Executing get_assets_by_category tool with category='{category}' and limit={limit}")
    data = await make_sams_request(f"/assets?category={category}&limit={limit}&page=1")

    if not data or "data" not in data:
        return f"Unable to fetch assets in category '{category}' from SAMS API."

    assets = data["data"]
    if not assets:
        return f"No assets found in category '{category}'."

    asset_list = []
    for asset in assets:
        asset_list.append(
            f"- ID: {asset.get('id')}\n"
            f"  Name: {asset.get('name', 'N/A')}\n"
            f"  Serial Number: {asset.get('serial_number', 'N/A')}\n"
            f"  Status: {asset.get('status', 'N/A')}\n"
            f"  Department: {asset.get('department', {}).get('name', 'N/A')}\n"
            f"  Current Value: ${asset.get('current_value', 0):,.2f}"
        )
    
    return f"Assets in Category '{category}' ({len(assets)} assets found):\n" + "\n".join(asset_list)

@mcp.tool()
async def get_assets_by_department(department: str, limit: int = 20) -> str:
    """
    Get assets filtered by their department.

    Args:
        department: The department name to filter by.
        limit: Maximum number of results to return. Defaults to 20.
    """
    logging.info(f"Executing get_assets_by_department tool with department='{department}' and limit={limit}")
    data = await make_sams_request(f"/assets?limit={limit}&page=1")

    if not data or "data" not in data:
        return f"Unable to fetch assets in department '{department}' from SAMS API."

    assets = data["data"]
    if not assets:
        return f"No assets found in department '{department}'."

    # Filter by department since the API doesn't support department filtering yet
    filtered_assets = [asset for asset in assets if asset.get('department', {}).get('name') == department]
    
    if not filtered_assets:
        return f"No assets found in department '{department}'."

    asset_list = []
    for asset in filtered_assets:
        asset_list.append(
            f"- ID: {asset.get('id')}\n"
            f"  Name: {asset.get('name', 'N/A')}\n"
            f"  Serial Number: {asset.get('serial_number', 'N/A')}\n"
            f"  Category: {asset.get('category', {}).get('name', 'N/A')}\n"
            f"  Status: {asset.get('status', 'N/A')}\n"
            f"  Current Value: ${asset.get('current_value', 0):,.2f}"
        )
    
    return f"Assets in Department '{department}' ({len(filtered_assets)} assets found):\n" + "\n".join(asset_list)

@mcp.tool()
async def get_asset_details(asset_id: str) -> str:
    """
    Get detailed information about a specific asset by its ID.

    Args:
        asset_id: The unique identifier of the asset.
    """
    logging.info(f"Executing get_asset_details tool with asset_id='{asset_id}'")
    data = await make_sams_request(f"/assets/{asset_id}")

    if not data or "data" not in data:
        return f"Unable to fetch asset details for ID '{asset_id}' from SAMS API."

    asset = data["data"]
    
    return (
        f"Asset Details for ID {asset_id}:\n"
        f"- Name: {asset.get('name', 'N/A')}\n"
        f"- Serial Number: {asset.get('serial_number', 'N/A')}\n"
        f"- Model: {asset.get('model', 'N/A')}\n"
        f"- Description: {asset.get('description', 'N/A')}\n"
        f"- Status: {asset.get('status', 'N/A')}\n"
        f"- Category: {asset.get('category', {}).get('name', 'N/A')}\n"
        f"- Department: {asset.get('department', {}).get('name', 'N/A')}\n"
        f"- Acquisition Date: {asset.get('acquisition_date', 'N/A')}\n"
        f"- Acquisition Cost: ${asset.get('acquisition_cost', 0):,.2f}\n"
        f"- Current Value: ${asset.get('current_value', 0):,.2f}\n"
        f"- Location: {asset.get('location', 'N/A')}\n"
        f"- Building: {asset.get('building', 'N/A')}\n"
        f"- Room: {asset.get('room', 'N/A')}\n"
        f"- Criticality: {asset.get('criticality', 'N/A')}\n"
        f"- Condition: {asset.get('condition', 'N/A')}"
    )

@mcp.tool()
async def get_assets_by_value_range(min_value: float = 0, max_value: float = 1000000, limit: int = 20) -> str:
    """
    Get assets within a specific value range.

    Args:
        min_value: Minimum asset value to include. Defaults to 0.
        max_value: Maximum asset value to include. Defaults to 1000000.
        limit: Maximum number of results to return. Defaults to 20.
    """
    logging.info(f"Executing get_assets_by_value_range tool with min_value={min_value}, max_value={max_value}, limit={limit}")
    data = await make_sams_request(f"/assets?limit={limit}&page=1")

    if not data or "data" not in data:
        return f"Unable to fetch assets from SAMS API for value range filtering."

    assets = data["data"]
    if not assets:
        return "No assets found."

    # Filter by value range since the API doesn't support value filtering yet
    filtered_assets = [
        asset for asset in assets 
        if asset.get('current_value') and min_value <= asset.get('current_value', 0) <= max_value
    ]
    
    if not filtered_assets:
        return f"No assets found with value between ${min_value:,.2f} and ${max_value:,.2f}."

    asset_list = []
    for asset in filtered_assets:
        asset_list.append(
            f"- ID: {asset.get('id')}\n"
            f"  Name: {asset.get('name', 'N/A')}\n"
            f"  Category: {asset.get('category', {}).get('name', 'N/A')}\n"
            f"  Department: {asset.get('department', {}).get('name', 'N/A')}\n"
            f"  Current Value: ${asset.get('current_value', 0):,.2f}\n"
            f"  Status: {asset.get('status', 'N/A')}"
        )
    
    return f"Assets with Value Between ${min_value:,.2f} and ${max_value:,.2f} ({len(filtered_assets)} assets found):\n" + "\n".join(asset_list)

@mcp.tool()
async def get_category_summary() -> str:
    """
    Get a summary of assets grouped by category with total values.
    """
    logging.info("Executing get_category_summary tool")
    data = await make_sams_request("/assets/summary-by-category")

    if not data or "data" not in data:
        return "Unable to fetch category summary from SAMS API."

    categories = data["data"]
    if not categories:
        return "No category data found."

    category_list = []
    for category in categories:
        category_list.append(
            f"- {category.get('category_name', 'Unknown')}: "
            f"{category.get('asset_count', 0)} assets, "
            f"Total Value: ${category.get('total_value', 0):,.2f}"
        )
    
    return "Assets by Category:\n" + "\n".join(category_list)

@mcp.tool()
async def get_status_summary() -> str:
    """
    Get a summary of assets grouped by status.
    """
    logging.info("Executing get_status_summary tool")
    data = await make_sams_request("/assets/summary-by-status")

    if not data or "data" not in data:
        return "Unable to fetch status summary from SAMS API."

    statuses = data["data"]
    if not statuses:
        return "No status data found."

    status_list = []
    for status in statuses:
        status_list.append(
            f"- {status.get('status', 'Unknown')}: "
            f"{status.get('asset_count', 0)} assets"
        )
    
    return "Assets by Status:\n" + "\n".join(status_list)

@mcp.tool()
async def get_assets_by_location(location_query: str, limit: int = 20) -> str:
    """
    Get assets filtered by location (address, building, room, or GPS coordinates).

    Args:
        location_query: The location to search for (e.g., "Jakarta", "Building A", "Room 101", "Main Office").
        limit: Maximum number of results to return. Defaults to 20.
    """
    logging.info(f"Executing get_assets_by_location tool with location_query='{location_query}' and limit={limit}")
    data = await make_sams_request(f"/assets?limit={limit}&page=1")

    if not data or "data" not in data:
        return f"Unable to fetch assets from SAMS API for location filtering."

    assets = data["data"]
    if not assets:
        return "No assets found."

    # Filter by location since the API doesn't support location filtering yet
    location_query_lower = location_query.lower()
    filtered_assets = []
    
    for asset in assets:
        # Check address, building_room, and coordinates
        address = asset.get('address', '').lower()
        building_room = asset.get('building_room', '').lower()
        
        # Check if location query matches any location field
        if (location_query_lower in address or 
            location_query_lower in building_room or
            (location_query_lower in "jakarta" and "jakarta" in address) or
            (location_query_lower in "office" and "office" in address) or
            (location_query_lower in "building" and "building" in building_room) or
            (location_query_lower in "room" and "room" in building_room)):
            filtered_assets.append(asset)
    
    if not filtered_assets:
        return f"No assets found at location '{location_query}'."

    asset_list = []
    for asset in filtered_assets:
        location_info = []
        if asset.get('address'):
            location_info.append(f"Address: {asset.get('address')}")
        if asset.get('building_room'):
            location_info.append(f"Building/Room: {asset.get('building_room')}")
        if asset.get('latitude') and asset.get('longitude'):
            location_info.append(f"Coordinates: {asset.get('latitude')}, {asset.get('longitude')}")
        
        location_str = " | ".join(location_info) if location_info else "Location: N/A"
        
        asset_list.append(
            f"- ID: {asset.get('id')}\n"
            f"  Name: {asset.get('name', 'N/A')}\n"
            f"  Serial Number: {asset.get('serial_number', 'N/A')}\n"
            f"  Category: {asset.get('category', {}).get('name', 'N/A')}\n"
            f"  Status: {asset.get('status', 'N/A')}\n"
            f"  Current Value: ${asset.get('current_value', 0):,.2f}\n"
            f"  {location_str}"
        )
    
    return f"Assets at Location '{location_query}' ({len(filtered_assets)} assets found):\n" + "\n".join(asset_list)

@mcp.tool()
async def get_assets_near_coordinates(latitude: float, longitude: float, radius_km: float = 1.0, limit: int = 20) -> str:
    """
    Get assets within a specified radius of GPS coordinates.

    Args:
        latitude: The latitude coordinate to search around.
        longitude: The longitude coordinate to search around.
        radius_km: Search radius in kilometers. Defaults to 1.0 km.
        limit: Maximum number of results to return. Defaults to 20.
    """
    logging.info(f"Executing get_assets_near_coordinates tool with lat={latitude}, lon={longitude}, radius={radius_km}km, limit={limit}")
    data = await make_sams_request(f"/assets?limit={limit}&page=1")

    if not data or "data" not in data:
        return f"Unable to fetch assets from SAMS API for coordinate-based filtering."

    assets = data["data"]
    if not assets:
        return "No assets found."

    # Filter assets by coordinate proximity
    filtered_assets = []
    for asset in assets:
        if asset.get('latitude') and asset.get('longitude'):
            # Calculate distance using Haversine formula (simplified)
            asset_lat = asset.get('latitude')
            asset_lon = asset.get('longitude')
            
            # Simple distance calculation (approximate)
            lat_diff = abs(asset_lat - latitude)
            lon_diff = abs(asset_lon - longitude)
            
            # Rough conversion to km (1 degree ≈ 111 km)
            distance_km = (lat_diff + lon_diff) * 111.0
            
            if distance_km <= radius_km:
                filtered_assets.append(asset)
    
    if not filtered_assets:
        return f"No assets found within {radius_km}km of coordinates ({latitude}, {longitude})."

    asset_list = []
    for asset in filtered_assets:
        asset_list.append(
            f"- ID: {asset.get('id')}\n"
            f"  Name: {asset.get('name', 'N/A')}\n"
            f"  Serial Number: {asset.get('serial_number', 'N/A')}\n"
            f"  Category: {asset.get('category', {}).get('name', 'N/A')}\n"
            f"  Status: {asset.get('status', 'N/A')}\n"
            f"  Current Value: ${asset.get('current_value', 0):,.2f}\n"
            f"  Coordinates: {asset.get('latitude')}, {asset.get('longitude')}\n"
            f"  Address: {asset.get('address', 'N/A')}"
        )
    
    return f"Assets within {radius_km}km of ({latitude}, {longitude}) ({len(filtered_assets)} assets found):\n" + "\n".join(asset_list)

# --- FastAPI HTTP Server ---

app = FastAPI()

@app.get("/")
async def root():
    """Root endpoint to test if the server is running."""
    return {"message": "SAMS MCP Server is running", "status": "healthy"}

@app.get("/tools")
async def list_tools():
    """Endpoint to get the list of available tools and their schemas."""
    # Get the registered tools from the FastMCP instance
    try:
        # Try to access tools through the mcp instance
        if hasattr(mcp, '_tools'):
            tools = mcp._tools
        elif hasattr(mcp, 'tools'):
            tools = mcp.tools
        else:
            # Fallback: manually list the tools we know about
            tools = {
                'get_asset_summary': get_asset_summary,
                'get_recent_assets': get_recent_assets,
                'search_assets': search_assets,
                'get_assets_by_status': get_assets_by_status,
                'get_assets_by_category': get_assets_by_category,
                'get_assets_by_department': get_assets_by_department,
                'get_asset_details': get_asset_details,
                'get_assets_by_value_range': get_assets_by_value_range,
                'get_category_summary': get_category_summary,
                'get_status_summary': get_status_summary,
                'get_assets_by_location': get_assets_by_location,
                'get_assets_near_coordinates': get_assets_near_coordinates
            }
        
        tool_schemas = []
        for name, tool_func in tools.items():
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
    except Exception as e:
        logging.error(f"Error listing tools: {e}")
        return {"error": str(e), "tools": []}

@app.post("/call/{tool_name}")
async def call_tool(tool_name: str, params: dict[str, Any]):
    """Endpoint to execute a specific tool with given parameters."""
    # Map tool names to functions
    tool_map = {
        'get_asset_summary': get_asset_summary,
        'get_recent_assets': get_recent_assets,
        'search_assets': search_assets,
        'get_assets_by_status': get_assets_by_status,
        'get_assets_by_category': get_assets_by_category,
        'get_assets_by_department': get_assets_by_department,
        'get_asset_details': get_asset_details,
        'get_assets_by_value_range': get_assets_by_value_range,
        'get_category_summary': get_category_summary,
        'get_status_summary': get_status_summary,
        'get_assets_by_location': get_assets_by_location,
        'get_assets_near_coordinates': get_assets_near_coordinates
    }
    
    if tool_name not in tool_map:
        raise HTTPException(status_code=404, detail=f"Tool '{tool_name}' not found.")

    tool_func = tool_map[tool_name]
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
