package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

// Structs for AI Chat
type AIChatRequest struct {
	Message string `json:"message"`
}

type AIChatResponse struct {
	Response string `json:"response"`
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
	log.Println("HandleAIChat: received new request")
	var req AIChatRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("HandleAIChat: error parsing request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("HandleAIChat: GEMINI_API_KEY not set")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "GEMINI_API_KEY is not set"})
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("HandleAIChat: failed to create genai client: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create AI client"})
	}
	defer client.Close()

	log.Println("HandleAIChat: fetching tools from MCP server")
	mcpTools, err := getToolsFromMCP()
	if err != nil {
		log.Printf("HandleAIChat: failed to get tools from MCP: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get tools from MCP server"})
	}
	log.Printf("HandleAIChat: successfully fetched %d tools from MCP server", len(mcpTools))

	model := client.GenerativeModel("gemini-1.5-pro-latest")
	if len(mcpTools) > 0 {
		model.Tools = mcpTools
	}

	cs := model.StartChat()

	log.Printf("HandleAIChat: sending prompt to Gemini: %s", req.Message)
	resp, err := cs.SendMessage(ctx, genai.Text(req.Message))
	if err != nil {
		log.Printf("HandleAIChat: failed to send message to Gemini: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get response from AI model"})
	}
	log.Println("HandleAIChat: successfully received response from Gemini")

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		log.Println("HandleAIChat: Gemini response has no content")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "AI model returned empty response"})
	}

	// Check if the model wants to call a tool
	if toolCalls := resp.Candidates[0].Content.Parts; len(toolCalls) > 0 {
		if tc, ok := toolCalls[0].(genai.ToolCall); ok {
			log.Printf("HandleAIChat: Gemini wants to call tool: %s", tc.Name)
			toolResult, err := callMCPTool(tc.Name, tc.Args)
			if err != nil {
				log.Printf("HandleAIChat: error calling MCP tool %s: %v", tc.Name, err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to execute tool"})
			}
			log.Printf("HandleAIChat: successfully called tool %s, got result", tc.Name)

			// Send the tool result back to the model
			resp, err = cs.SendMessage(ctx, genai.ToolResponse{
				Name:    tc.Name,
				Content: toolResult,
			})
			if err != nil {
				log.Printf("HandleAIChat: error sending tool response to Gemini: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send tool response to AI model"})
			}
			log.Println("HandleAIChat: successfully sent tool response to Gemini and got final answer")
		}
	}

	// Extract and send the final text response
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if finalResponse, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			log.Printf("HandleAIChat: sending final response to client: %s", finalResponse)
			return c.JSON(AIChatResponse{Response: string(finalResponse)})
		}
	}

	log.Println("HandleAIChat: could not extract final text response from Gemini")
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not extract final response from AI model"})
}
