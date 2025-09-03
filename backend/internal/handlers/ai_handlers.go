package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Constants
const (
	mcpServerBaseURL = "http://sams-mcp-server:8081"
)

// Structs for Gemini API
type GeminiRequest struct {
	Prompt string `json:"prompt"`
}

// Struct for MCP Tool Schemas
type MCPToolProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type MCPToolParameters struct {
	Type       string                     `json:"type"`
	Properties map[string]MCPToolProperty `json:"properties"`
	Required   []string                   `json:"required"`
}

type MCPTool struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Parameters  MCPToolParameters `json:"parameters"`
}

type MCPToolsResponse struct {
	Tools []MCPTool `json:"tools"`
}

// getToolsFromMCP fetches tool definitions from the MCP server.
func getToolsFromMCP() ([]*genai.Tool, error) {
	resp, err := http.Get(mcpServerBaseURL + "/tools")
	if err != nil {
		return nil, fmt.Errorf("failed to contact MCP server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MCP server returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read MCP server response: %w", err)
	}

	var mcpToolsResp MCPToolsResponse
	if err := json.Unmarshal(body, &mcpToolsResp); err != nil {
		return nil, fmt.Errorf("failed to parse MCP tools response: %w", err)
	}

	// Convert MCP tools to Gemini's genai.Tool format
	var geminiTools []*genai.Tool
	functionDeclarations := []*genai.FunctionDeclaration{}

	for _, t := range mcpToolsResp.Tools {
		props := make(map[string]*genai.Schema)
		for name, p := range t.Parameters.Properties {
			var genaiType genai.Type
			switch p.Type {
			case "integer":
				genaiType = genai.TypeInteger
			case "number":
				genaiType = genai.TypeNumber
			case "boolean":
				genaiType = genai.TypeBoolean
			default:
				genaiType = genai.TypeString
			}
			props[name] = &genai.Schema{Type: genaiType, Description: p.Description}
		}

		fd := &genai.FunctionDeclaration{
			Name:        t.Name,
			Description: t.Description,
			Parameters: &genai.Schema{
				Type:       genai.TypeObject,
				Properties: props,
				Required:   t.Parameters.Required,
			},
		}
		functionDeclarations = append(functionDeclarations, fd)
	}

	geminiTools = append(geminiTools, &genai.Tool{FunctionDeclarations: functionDeclarations})
	return geminiTools, nil
}

// callMCPTool calls a specific tool on the MCP server.
func callMCPTool(toolName string, params map[string]any) (map[string]any, error) {
	requestBody, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tool params: %w", err)
	}

	resp, err := http.Post(mcpServerBaseURL+"/call/"+toolName, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to call tool on MCP server: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read tool response body: %w", err)
	}

	// The MCP server returns {"result": "some string"}. We need to extract the string
	// and package it correctly for Gemini.
	var toolCallResult struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(body, &toolCallResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tool result from MCP: %w", err)
	}

	// Gemini expects a map[string]any as the result, so we wrap it.
	resultForGemini := map[string]any{"output": toolCallResult.Result}
	return resultForGemini, nil
}

// HandleAIChat is the main handler for the AI chat functionality.
func HandleAIChat(c *fiber.Ctx) error {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "GEMINI_API_KEY is not set"})
	}

	var req GeminiRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	ctx := c.Context()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create Gemini client", "details": err.Error()})
	}
	defer client.Close()

	tools, err := getToolsFromMCP()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get tools from MCP server", "details": err.Error()})
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	model.Tools = tools

	session := model.StartChat()
	resp, err := session.SendMessage(ctx, genai.Text(req.Prompt))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send message to Gemini", "details": err.Error()})
	}

	// Handle tool calls if any
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if fc := resp.Candidates[0].Content.Parts[0].(genai.FunctionCall); fc.Name != "" {

			toolResult, err := callMCPTool(fc.Name, fc.Args)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to execute tool", "details": err.Error()})
			}

			// Send the tool result back to Gemini
			resp, err = session.SendMessage(ctx, genai.FunctionResponse{
				Name:     fc.Name,
				Response: toolResult,
			})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send tool response to Gemini", "details": err.Error()})
			}
		}
	}

	// Extract and send the final text response
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return c.JSON(fiber.Map{"response": text})
		}
	}

	return c.JSON(fiber.Map{"response": "No content received from AI."})
}
