package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"sams-backend/internal/database"
	"sams-backend/internal/models"
)

// AIQueryRequest represents the request body for AI queries
type AIQueryRequest struct {
	Query string `json:"query"`
}

// GeminiRequest represents the request to Gemini API
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// GeminiContent represents content for Gemini API
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a part of content for Gemini API
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse represents the response from Gemini API
type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

// GeminiCandidate represents a candidate response from Gemini API
type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

func HandleAIQuery(c *fiber.Ctx) error {
	db := database.GetDB()
	var request AIQueryRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid request body"})
	}

	if request.Query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Query is required"})
	}

	// 1. Get relevant assets from the database
	assets, err := getRelevantAssets(db, request.Query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to retrieve asset context"})
	}

	// 2. Prepare context for Gemini
	context := prepareAssetContext(assets)

	// 3. Call Gemini API
	response, err := callGeminiAPI(request.Query, context)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"error": false, "data": response})
}

// getRelevantAssets fetches assets relevant to the query
func getRelevantAssets(db *gorm.DB, query string) ([]models.AssetSummary, error) {
	var assets []models.AssetSummary

	// Simple keyword-based filtering for MVP
	query = strings.ToLower(query)

	baseQuery := db.Table("asset_summary")

	if strings.Contains(query, "it") || strings.Contains(query, "computer") || strings.Contains(query, "laptop") {
		baseQuery = baseQuery.Where("category_name ILIKE ?", "%IT%")
	} else if strings.Contains(query, "vehicle") || strings.Contains(query, "car") || strings.Contains(query, "van") {
		baseQuery = baseQuery.Where("category_name ILIKE ?", "%Vehicle%")
	} else if strings.Contains(query, "building") || strings.Contains(query, "office") {
		baseQuery = baseQuery.Where("category_name ILIKE ?", "%Building%")
	} else if strings.Contains(query, "machinery") || strings.Contains(query, "equipment") {
		baseQuery = baseQuery.Where("category_name ILIKE ?", "%Machinery%")
	}

	if strings.Contains(query, "active") || strings.Contains(query, "working") {
		baseQuery = baseQuery.Where("status = ?", "active")
	} else if strings.Contains(query, "maintenance") {
		baseQuery = baseQuery.Where("status = ?", "maintenance")
	}

	if strings.Contains(query, "high") || strings.Contains(query, "critical") {
		baseQuery = baseQuery.Where("criticality IN ?", []string{"high", "critical"})
	}

	// Limit results for context
	if err := baseQuery.Limit(10).Find(&assets).Error; err != nil {
		return nil, err
	}

	// If no specific filtering, get a sample of assets
	if len(assets) == 0 {
		if err := db.Table("asset_summary").Limit(5).Find(&assets).Error; err != nil {
			return nil, err
		}
	}

	return assets, nil
}

// prepareAssetContext prepares asset data context for AI
func prepareAssetContext(assets []models.AssetSummary) string {
	if len(assets) == 0 {
		return "No specific asset data available for this query."
	}

	var context strings.Builder
	context.WriteString(fmt.Sprintf("Found %d relevant assets:\n\n", len(assets)))

	for i, asset := range assets {
		context.WriteString(fmt.Sprintf("Asset %d:\n", i+1))
		context.WriteString(fmt.Sprintf("- Name: %s\n", asset.Name))
		context.WriteString(fmt.Sprintf("- Category: %s\n", asset.CategoryName))
		context.WriteString(fmt.Sprintf("- Type: %s\n", asset.Type))
		context.WriteString(fmt.Sprintf("- Status: %s\n", asset.Status))
		context.WriteString(fmt.Sprintf("- Condition: %s\n", asset.Condition))
		context.WriteString(fmt.Sprintf("- Criticality: %s\n", asset.Criticality))
		context.WriteString(fmt.Sprintf("- Location: %s\n", asset.Address))
		if asset.Latitude != nil && asset.Longitude != nil {
			context.WriteString(fmt.Sprintf("- Coordinates: %.6f, %.6f\n", *asset.Latitude, *asset.Longitude))
		}
		context.WriteString(fmt.Sprintf("- Current Value: $%.2f\n", asset.CurrentValue))
		context.WriteString("\n")
	}

	return context.String()
}

// callGeminiAPI calls the Gemini API with the given prompt
func callGeminiAPI(apiKey, prompt string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", apiKey)

	requestBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Gemini API returned status: %d", resp.StatusCode)
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "No response generated from AI.", nil
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
