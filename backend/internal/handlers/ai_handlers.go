package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Structs for AI Chat
type AIChatRequest struct {
	Message string `json:"message"`
}

type AIChatResponse struct {
	Response string `json:"response"`
}

// MCP Server Integration
type MCPToolCall struct {
	ToolName string                 `json:"tool_name"`
	Params   map[string]interface{} `json:"params"`
}

type MCPToolResponse struct {
	Result string `json:"result"`
}

// Asset-related keywords to detect when to use MCP tools
var assetKeywords = []string{
	"asset", "assets", "laptop", "computer", "equipment", "vehicle", "camera",
	"printer", "tablet", "phone", "it", "hardware", "inventory", "summary",
	"total", "value", "status", "category", "department", "search", "find",
	"show", "list", "count", "maintenance", "active", "inactive", "disposed",
}

// Check if the message is asset-related
func isAssetRelated(message string) bool {
	lowerMessage := strings.ToLower(message)

	// Check for generic asset keywords
	for _, keyword := range assetKeywords {
		if strings.Contains(lowerMessage, keyword) {
			return true
		}
	}

	// Also check for specific asset references (brands, models, etc.)
	if containsSpecificAssetReference(message) {
		return true
	}

	return false
}

// Check if a word is a common word that doesn't indicate a specific asset
func isCommonWord(word string) bool {
	commonWords := []string{
		"the", "a", "an", "and", "or", "but", "in", "on", "at", "to", "for", "of", "with",
		"by", "from", "up", "about", "into", "through", "during", "before", "after", "above",
		"below", "between", "among", "asset", "assets", "show", "give", "me", "info", "about",
		"what", "is", "are", "where", "when", "how", "why", "which", "who", "this", "that",
		"these", "those", "here", "there", "now", "then", "today", "yesterday", "tomorrow",
	}

	lowerWord := strings.ToLower(word)
	for _, common := range commonWords {
		if lowerWord == common {
			return true
		}
	}

	return false
}

// Check if the message contains a specific asset reference
func containsSpecificAssetReference(message string) bool {
	// Common asset brands and models that indicate specific asset queries
	specificAssetKeywords := []string{
		"samsung", "galaxy", "tab", "lenovo", "thinkpad", "dell", "latitude", "hp", "elitebook",
		"canon", "eos", "nikon", "sony", "apple", "macbook", "iphone", "ipad", "toyota", "honda",
		"bosch", "lg", "epson", "cisco", "microsoft", "surface", "yamaha", "psr",
	}

	lowerMessage := strings.ToLower(message)

	// Check if message contains specific asset keywords (brands/models)
	for _, keyword := range specificAssetKeywords {
		if strings.Contains(lowerMessage, keyword) {
			return true
		}
	}

	// Check for specific asset patterns (e.g., "Samsung Galaxy Tab S7", "Lenovo ThinkPad X1")
	// Only if the message looks like it's asking about a specific asset, not general questions
	words := strings.Fields(message)

	// Skip if this looks like a general question (starts with question words)
	questionWords := []string{"what", "how", "when", "where", "why", "which", "who", "show", "give", "tell", "find", "search", "get"}
	if len(words) > 0 {
		firstWord := strings.ToLower(strings.Trim(words[0], "?!.,"))
		for _, qw := range questionWords {
			if firstWord == qw {
				// This is a general question, not a specific asset reference
				return false
			}
		}
	}

	// Check if message contains what looks like a specific asset name
	// Only if it doesn't contain general query words
	generalQueryWords := []string{"total", "value", "cost", "worth", "summary", "overview", "count", "how many", "equipment", "assets", "category"}
	for _, gqw := range generalQueryWords {
		if strings.Contains(lowerMessage, gqw) {
			// This is a general query, not a specific asset reference
			return false
		}
	}

	// If we get here, check if it looks like a specific asset name
	if len(words) >= 2 {
		// Look for brand + model patterns
		hasBrand := false
		hasModel := false

		for _, word := range words {
			lowerWord := strings.ToLower(word)
			if len(word) > 2 && !isCommonWord(lowerWord) {
				// Check if this could be a brand or model
				if contains(specificAssetKeywords, lowerWord) {
					hasBrand = true
				} else if len(word) >= 3 && !isCommonWord(lowerWord) {
					hasModel = true
				}
			}
		}

		// Only return true if we have both brand and model indicators
		return hasBrand && hasModel
	}

	return false
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Extract the actual asset name from a query message
func extractAssetName(message string) string {
	// Remove common query words and phrases
	queryWords := []string{
		"give me info about asset", "show me", "tell me about", "what is", "find", "search for",
		"get info about", "information about", "details about", "show", "give", "me", "info", "about",
	}

	cleanedMessage := strings.ToLower(message)
	for _, word := range queryWords {
		cleanedMessage = strings.ReplaceAll(cleanedMessage, strings.ToLower(word), "")
	}

	// Clean up extra spaces and trim
	cleanedMessage = strings.TrimSpace(cleanedMessage)
	cleanedMessage = strings.Join(strings.Fields(cleanedMessage), " ")

	// If we have a meaningful asset name, return it
	if len(cleanedMessage) > 0 && len(cleanedMessage) < 100 {
		return cleanedMessage
	}

	// Fallback: try to extract brand/model combinations
	words := strings.Fields(message)
	assetWords := []string{}

	for _, word := range words {
		lowerWord := strings.ToLower(word)
		if !isCommonWord(lowerWord) && len(word) > 2 {
			assetWords = append(assetWords, word)
		}
	}

	if len(assetWords) > 0 {
		return strings.Join(assetWords, " ")
	}

	// Last resort: return the original message
	return message
}

// Call MCP server tool
func callMCPTool(toolName string, params map[string]interface{}) (string, error) {
	mcpURL := "http://sams-mcp-server:8081/call/" + toolName

	// Convert params to JSON
	jsonData, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("failed to marshal params: %v", err)
	}

	// Make request to MCP server
	resp, err := http.Post(mcpURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to call MCP server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("MCP server returned status: %d", resp.StatusCode)
	}

	// Parse response
	var mcpResp MCPToolResponse
	if err := json.NewDecoder(resp.Body).Decode(&mcpResp); err != nil {
		return "", fmt.Errorf("failed to decode MCP response: %v", err)
	}

	return mcpResp.Result, nil
}

// Determine which MCP tool to use based on the message
func determineMCPTool(message string) (string, map[string]interface{}) {
	lowerMessage := strings.ToLower(message)

	// Check for specific asset queries first (before generic patterns)
	// If the message contains specific asset names, models, or serial numbers, use search
	if containsSpecificAssetReference(message) {
		assetName := extractAssetName(message)
		return "search_assets", map[string]interface{}{"query": assetName, "limit": 10}
	}

	// PRIORITY 1: Value and summary queries (highest priority)
	if strings.Contains(lowerMessage, "total value") || strings.Contains(lowerMessage, "total cost") ||
		strings.Contains(lowerMessage, "total worth") || strings.Contains(lowerMessage, "value of") ||
		strings.Contains(lowerMessage, "cost of") || strings.Contains(lowerMessage, "worth of") {

		// If asking about specific category value, use category tool
		if strings.Contains(lowerMessage, "it") || strings.Contains(lowerMessage, "it equipment") {
			return "get_assets_by_category", map[string]interface{}{"category": "IT Equipment", "limit": 50}
		}
		if strings.Contains(lowerMessage, "vehicle") || strings.Contains(lowerMessage, "vehicles") {
			return "get_assets_by_category", map[string]interface{}{"category": "Vehicles", "limit": 50}
		}
		if strings.Contains(lowerMessage, "tool") || strings.Contains(lowerMessage, "tools") {
			return "get_assets_by_category", map[string]interface{}{"category": "Tools", "limit": 50}
		}

		// For general value queries, get all assets to calculate total
		return "get_asset_summary", map[string]interface{}{}
	}

	// PRIORITY 2: Category-specific queries
	if strings.Contains(lowerMessage, "it equipment") || (strings.Contains(lowerMessage, "it") && strings.Contains(lowerMessage, "equipment")) {
		return "get_assets_by_category", map[string]interface{}{"category": "IT Equipment", "limit": 20}
	}
	if strings.Contains(lowerMessage, "vehicle") || strings.Contains(lowerMessage, "vehicles") {
		return "get_assets_by_category", map[string]interface{}{"category": "Vehicles", "limit": 20}
	}
	if strings.Contains(lowerMessage, "tool") || strings.Contains(lowerMessage, "tools") {
		return "get_assets_by_category", map[string]interface{}{"category": "Tools", "limit": 20}
	}

	// PRIORITY 2.5: Department-specific queries
	if strings.Contains(lowerMessage, "department") || strings.Contains(lowerMessage, "departement") ||
		strings.Contains(lowerMessage, "dept") || strings.Contains(lowerMessage, "by department") {

		// First, try to extract the actual department name from the query
		departmentName := ""

		// Check for specific department names from our database
		if strings.Contains(lowerMessage, "project") {
			departmentName = "Project"
		} else if strings.Contains(lowerMessage, "finance") || strings.Contains(lowerMessage, "accounting") {
			departmentName = "Finance"
		} else if strings.Contains(lowerMessage, "human capital") || strings.Contains(lowerMessage, "hc") {
			departmentName = "Human Capital (HC)"
		} else if strings.Contains(lowerMessage, "operation") || strings.Contains(lowerMessage, "operations") || strings.Contains(lowerMessage, "ops") {
			departmentName = "Operation"
		} else if strings.Contains(lowerMessage, "information technology") || strings.Contains(lowerMessage, "it") {
			departmentName = "Information Technology (IT)"
		} else if strings.Contains(lowerMessage, "marketing") {
			departmentName = "Marketing"
		}

		// If no specific department found, default to IT
		if departmentName == "" {
			departmentName = "Information Technology (IT)"
		}

		return "get_assets_by_department", map[string]interface{}{"department": departmentName, "limit": 20}
	}

	// PRIORITY 3: Asset summary and overview queries
	if strings.Contains(lowerMessage, "summary") || strings.Contains(lowerMessage, "overview") ||
		strings.Contains(lowerMessage, "count") || strings.Contains(lowerMessage, "how many") {
		return "get_asset_summary", map[string]interface{}{}
	}

	// PRIORITY 4: Status-based queries
	if strings.Contains(lowerMessage, "active") {
		return "get_assets_by_status", map[string]interface{}{"status": "active", "limit": 20}
	}
	if strings.Contains(lowerMessage, "maintenance") {
		return "get_assets_by_status", map[string]interface{}{"status": "maintenance", "limit": 20}
	}

	// PRIORITY 5: Location-based queries
	if strings.Contains(lowerMessage, "location") || strings.Contains(lowerMessage, "address") ||
		strings.Contains(lowerMessage, "building") || strings.Contains(lowerMessage, "room") ||
		strings.Contains(lowerMessage, "jakarta") || strings.Contains(lowerMessage, "office") ||
		strings.Contains(lowerMessage, "where") || strings.Contains(lowerMessage, "place") {
		// Extract location query from the message
		locationQuery := message
		if strings.Contains(lowerMessage, "jakarta") {
			locationQuery = "Jakarta"
		} else if strings.Contains(lowerMessage, "office") {
			locationQuery = "Office"
		} else if strings.Contains(lowerMessage, "building") {
			locationQuery = "Building"
		} else if strings.Contains(lowerMessage, "room") {
			locationQuery = "Room"
		}
		return "get_assets_by_location", map[string]interface{}{"location_query": locationQuery, "limit": 20}
	}

	// PRIORITY 6: Generic search queries (lowest priority)
	if strings.Contains(lowerMessage, "search") || strings.Contains(lowerMessage, "find") ||
		strings.Contains(lowerMessage, "computer") {
		// Extract search query
		query := message
		if strings.Contains(lowerMessage, "laptop") {
			query = "laptop"
		} else if strings.Contains(lowerMessage, "computer") {
			query = "computer"
		}
		return "search_assets", map[string]interface{}{"query": query, "limit": 10}
	}

	// Default to asset summary
	return "get_asset_summary", map[string]interface{}{}
}

// HandleAIQuery is the main handler for the AI chat functionality.
func HandleAIQuery(c *fiber.Ctx) error {
	log.Println("HandleAIQuery: received new request")
	var req AIChatRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("HandleAIQuery: error parsing request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	// FORCE GEMINI TO USE TOOLS: Always use Gemini but with tool knowledge
	log.Printf("HandleAIQuery: FORCE GEMINI TOOL MODE - instructing Gemini to use MCP tools for: %s", req.Message)

	// Determine which MCP tool to use
	toolName, params := determineMCPTool(req.Message)

	// Call MCP server to get actual data
	mcpResult, err := callMCPTool(toolName, params)
	if err != nil {
		log.Printf("HandleAIQuery: MCP tool call failed: %v", err)
		// Continue to Gemini with tool information even if MCP fails
		mcpResult = "Tool call failed: " + err.Error()
	}

	// Now send to Gemini with tool data and instructions to use it
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	log.Printf("HandleAIQuery: GEMINI_API_KEY value: %s", apiKey)
	if apiKey == "" || apiKey == "your-google-ai-api-key-here" {
		log.Println("HandleAIQuery: GEMINI_API_KEY not set or is placeholder, providing fallback response")
		// Provide a fallback response with the MCP tool data
		fallbackResponse := fmt.Sprintf(`Hello! I'm your SAMS AI Assistant. I can help you with asset management questions.

Your question: "%s"

Here's what I found:
%s

I can help you with questions about your assets, such as:
• Asset summaries and total values (in Indonesian Rupiah)
• Finding assets by category, department, or location
• Asset status and maintenance information
• Search for specific assets

What would you like to know about your assets? Remember, all monetary values are presented in Indonesian Rupiah (Rp).`, req.Message, mcpResult)

		return c.JSON(fiber.Map{
			"response": fallbackResponse,
		})
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("HandleAIQuery: failed to create genai client: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create AI client"})
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")
	cs := model.StartChat()

	// Enhanced prompt that forces Gemini to use tool data
	enhancedPrompt := fmt.Sprintf(`You are a helpful SAMS (Smart Asset Management System) AI Assistant. You help users with asset management questions by providing accurate information from the system.

User Query: %s

Here is the current data from the system:
%s

INSTRUCTIONS:
1. Answer the user's question naturally and conversationally using ONLY the data provided above
2. Do NOT mention any technical details about tools, APIs, or data retrieval methods
3. Respond as if you naturally know this information about their assets
4. If the data is empty or shows an error, politely explain that you couldn't find the information and suggest they try rephrasing their question
5. Be helpful and friendly in your responses
6. Format ALL currency values in Indonesian Rupiah (Rp) format with commas for thousands (e.g., "Rp5,000,000.00")
7. If asked about specific assets, categories, or departments, provide the relevant information from the data

Remember: You are a knowledgeable assistant who helps with asset management. Respond naturally without mentioning how you got the information. Always present monetary values in Indonesian Rupiah (Rp).`,
		req.Message, mcpResult)

	log.Printf("HandleAIQuery: sending tool-enhanced prompt to Gemini: %s", enhancedPrompt)
	resp, err := cs.SendMessage(ctx, genai.Text(enhancedPrompt))
	if err != nil {
		log.Printf("HandleAIQuery: failed to send message to Gemini: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get response from AI model"})
	}
	log.Println("HandleAIQuery: successfully received response from Gemini")

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		log.Println("HandleAIQuery: Gemini response has no content")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "AI model returned empty response"})
	}

	// Extract and send the final text response
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if finalResponse, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			log.Printf("HandleAIQuery: sending final Gemini response to client: %s", string(finalResponse))
			return c.JSON(AIChatResponse{Response: string(finalResponse)})
		}
	}

	log.Println("HandleAIQuery: could not extract final text response from Gemini")
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not extract final response from AI model"})
}

// Generate a helpful response when tools fail, encouraging tool usage
func generateToolFocusedResponse(message string, attemptedTool string, err error) string {
	// Provide helpful guidance based on the attempted tool
	switch attemptedTool {
	case "search_assets":
		return `I tried to search for assets but encountered an issue. Here are some ways to query your assets:

**Available Asset Queries:**
• Search by name: "Samsung Galaxy Tab S7", "Lenovo ThinkPad X1"
• Search by category: "IT Equipment", "Vehicles", "Tools"
• Search by status: "Active assets", "Maintenance assets"
• Search by location: "Assets in Jakarta", "Office assets"
• Get summaries: "Asset summary", "Total value"

**Try rephrasing your query** to be more specific about what you want to know about your assets.`

	case "get_asset_summary":
		return `I tried to get asset summary information but encountered an issue. 

**Available Summary Queries:**
• "Show me total assets and value"
• "Asset overview"
• "How many assets do we have?"
• "Total asset count"

**Try asking for specific asset information** or rephrase your summary request.`

	case "get_assets_by_location":
		return `I tried to find assets by location but encountered an issue.

**Available Location Queries:**
• "Assets in Jakarta"
• "Show me office assets"
• "Building assets"
• "Room 101 assets"

**Try being more specific** about the location you're interested in.`

	default:
		return `I tried to use asset management tools but encountered an issue. 

**What I can help you with:**
• **Asset Information**: Search for specific assets by name, model, or serial number
• **Asset Categories**: Find assets by type (IT Equipment, Vehicles, Tools)
• **Asset Status**: Find active, maintenance, or disposed assets
• **Asset Location**: Find assets by location, building, or room
• **Asset Summaries**: Get total counts, values, and overviews

**Try asking about your assets** in a specific way, such as:
• "Show me laptop assets"
• "Find Canon camera"
• "Assets in Jakarta"
• "Total asset value"`

	}
}
