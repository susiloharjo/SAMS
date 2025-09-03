package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"

	"sams-backend/internal/database"
	"sams-backend/internal/models"
)

type CategorySummary struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// GetCategorySummary godoc
// @Summary Get asset summary by category
// @Description Get a summary of asset values grouped by category
// @Tags assets
// @Accept  json
// @Produce  json
// @Success 200 {array} CategorySummary
// @Failure 500 {object} fiber.Map
// @Router /assets/summary-by-category [get]
func GetCategorySummary(c *fiber.Ctx) error {
	db := database.GetDB()
	var results []CategorySummary

	err := db.Table("assets").
		Select("categories.name, SUM(assets.current_value) as value").
		Joins("LEFT JOIN categories ON categories.id = assets.category_id").
		Group("categories.name").
		Having("SUM(assets.current_value) > 0").
		Scan(&results).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to get category summary"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "data": results})
}

type AssetSummary struct {
	TotalAssets    int64   `json:"total_assets"`
	TotalValue     float64 `json:"total_value"`
	ActiveAssets   int64   `json:"active_assets"`
	CriticalAssets int64   `json:"critical_assets"`
}

// GetAssetSummary godoc
// @Summary Get asset summary statistics
// @Description Get summary statistics for all assets
// @Tags assets
// @Accept  json
// @Produce  json
// @Success 200 {object} AssetSummary
// @Failure 500 {object} fiber.Map
// @Router /assets/summary [get]
func GetAssetSummary(c *fiber.Ctx) error {
	db := database.GetDB()
	var summary AssetSummary

	// Get total assets
	if err := db.Model(&models.Asset{}).Count(&summary.TotalAssets).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to get total assets count"})
	}

	// Get total value
	if err := db.Model(&models.Asset{}).Select("COALESCE(sum(current_value), 0)").Row().Scan(&summary.TotalValue); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to get total asset value"})
	}

	// Get active assets
	if err := db.Model(&models.Asset{}).Where("status = ?", "active").Count(&summary.ActiveAssets).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to get active assets count"})
	}

	// Get critical assets
	if err := db.Model(&models.Asset{}).Where("criticality = ?", "critical").Count(&summary.CriticalAssets).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to get critical assets count"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "data": summary})
}

// GetAssets godoc
// @Summary Get all assets
// @Description Get a paginated list of all assets
// @Tags assets
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search term"
// @Param category query string false "Category name"
// @Param status query string false "Asset status"
// @Success 200 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /assets [get]
func GetAssets(c *fiber.Ctx) error {
	db := database.GetDB()
	var assets []models.Asset
	var total int64

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// Ensure valid pagination values
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Search functionality
	if search := c.Query("search"); search != "" {
		searchTerm := "%" + search + "%"
		db = db.Where("name ILIKE ? OR serial_number ILIKE ? OR model ILIKE ? OR description ILIKE ?", searchTerm, searchTerm, searchTerm, searchTerm)
	}

	// Apply filters
	if categoryName := c.Query("category"); categoryName != "" && categoryName != "all" {
		// Find category ID by name
		var category models.Category
		if err := db.Where("name = ?", categoryName).First(&category).Error; err == nil {
			db = db.Where("category_id = ?", category.ID)
		}
	}

	if status := c.Query("status"); status != "" && status != "all" {
		db = db.Where("status = ?", status)
	}

	if condition := c.Query("condition"); condition != "" {
		db = db.Where("condition = ?", condition)
	}

	if err := db.Model(&models.Asset{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to count assets",
		})
	}

	if err := db.Preload("Category").Preload("Department").Offset(offset).Limit(limit).Order("created_at DESC").Find(&assets).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch assets",
		})
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return c.JSON(fiber.Map{
		"error":   false,
		"data":    assets,
		"message": "Assets retrieved successfully",
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetAsset returns a single asset by ID
// @Failure 500 {object} fiber.Map
// @Router /assets/{id} [get]
func GetAsset(c *fiber.Ctx) error {
	db := database.GetDB()
	var asset models.Asset
	id := c.Params("id")
	assetID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid asset ID",
		})
	}

	if err := db.Preload("Category").Preload("Department").First(&asset, "id = ?", assetID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Asset not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch asset",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"data":    asset,
		"message": "Asset retrieved successfully",
	})
}

// CreateAsset creates a new asset
// @Failure 500 {object} fiber.Map
// @Router /assets [post]
func CreateAsset(c *fiber.Ctx) error {
	db := database.GetDB()
	var asset models.Asset

	if err := c.BodyParser(&asset); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	if asset.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Asset name is required",
		})
	}

	if asset.SerialNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Serial number is required",
		})
	}

	// Check if serial number already exists
	var existingAsset models.Asset
	if err := db.Where("serial_number = ?", asset.SerialNumber).First(&existingAsset).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "Asset with this serial number already exists",
		})
	}

	// Set default values
	if asset.Status == "" {
		asset.Status = "active"
	}
	if asset.Condition == "" {
		asset.Condition = "good"
	}
	if asset.Criticality == "" {
		asset.Criticality = "low"
	}

	if err := db.Create(&asset).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create asset",
		})
	}

	// Reload with category information
	db.Preload("Category").First(&asset, asset.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"data":    asset,
		"message": "Asset created successfully",
	})
}

// UpdateAsset updates an existing asset
// @Failure 500 {object} fiber.Map
// @Router /assets/{id} [put]
func UpdateAsset(c *fiber.Ctx) error {
	db := database.GetDB()
	var asset models.Asset
	id := c.Params("id")
	assetID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid asset ID",
		})
	}

	if err := db.First(&asset, "id = ?", assetID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Asset not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch asset",
		})
	}

	var updateData models.Asset
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Update fields if provided
	if updateData.Name != "" {
		asset.Name = updateData.Name
	}
	if updateData.Description != "" {
		asset.Description = updateData.Description
	}
	if updateData.CategoryID != nil {
		asset.CategoryID = updateData.CategoryID
	}
	if updateData.DepartmentID != nil {
		asset.DepartmentID = updateData.DepartmentID
	}
	if updateData.Type != "" {
		asset.Type = updateData.Type
	}
	if updateData.Model != "" {
		asset.Model = updateData.Model
	}
	if updateData.Manufacturer != "" {
		asset.Manufacturer = updateData.Manufacturer
	}
	if updateData.AcquisitionCost != 0 {
		asset.AcquisitionCost = updateData.AcquisitionCost
	}
	if updateData.CurrentValue != 0 {
		asset.CurrentValue = updateData.CurrentValue
	}
	if updateData.DepreciationRate != 0 {
		asset.DepreciationRate = updateData.DepreciationRate
	}
	if updateData.Status != "" {
		asset.Status = updateData.Status
	}
	if updateData.Condition != "" {
		asset.Condition = updateData.Condition
	}
	if updateData.Criticality != "" {
		asset.Criticality = updateData.Criticality
	}
	if updateData.Latitude != nil {
		asset.Latitude = updateData.Latitude
	}
	if updateData.Longitude != nil {
		asset.Longitude = updateData.Longitude
	}
	if updateData.Address != "" {
		asset.Address = updateData.Address
	}
	if updateData.BuildingRoom != "" {
		asset.BuildingRoom = updateData.BuildingRoom
	}
	if updateData.AcquisitionDate != nil {
		asset.AcquisitionDate = updateData.AcquisitionDate
	}
	if updateData.ExpectedLifeYears != nil {
		asset.ExpectedLifeYears = updateData.ExpectedLifeYears
	}
	if updateData.MaintenanceSchedule != "" {
		asset.MaintenanceSchedule = updateData.MaintenanceSchedule
	}
	if updateData.Certifications != "" {
		asset.Certifications = updateData.Certifications
	}
	if updateData.Standards != "" {
		asset.Standards = updateData.Standards
	}
	if updateData.AuditInfo != "" {
		asset.AuditInfo = updateData.AuditInfo
	}

	if err := db.Save(&asset).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update asset",
		})
	}

	// Reload with category and department information
	db.Preload("Category").Preload("Department").First(&asset, "id = ?", id)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"data":    asset,
		"message": "Asset updated successfully",
	})
}

// DeleteAsset deletes an asset
// @Failure 500 {object} fiber.Map
// @Router /assets/{id} [delete]
func DeleteAsset(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")
	assetID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid asset ID",
		})
	}

	var asset models.Asset
	if err := db.First(&asset, "id = ?", assetID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Asset not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch asset",
		})
	}

	if err := db.Delete(&asset).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete asset",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Asset deleted successfully",
	})
}

// GenerateAssetQR generates a QR code for an asset
// @Failure 500 {object} fiber.Map
// @Router /assets/{id}/qr [get]
func GenerateAssetQR(c *fiber.Ctx) error {
	db := database.GetDB()
	var asset models.Asset
	id := c.Params("id")
	assetID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid asset ID",
		})
	}

	if err := db.First(&asset, "id = ?", assetID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Asset not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch asset",
		})
	}

	// Create QR code content
	qrContent := fmt.Sprintf("Asset: %s\nSerial: %s\nCategory: %s\nLocation: %s",
		asset.Name, asset.SerialNumber, asset.Type, asset.Address)

	// Generate QR code
	qrCode, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to generate QR code",
		})
	}

	// Set response headers
	c.Set("Content-Type", "image/png")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=asset_%s_qr.png", asset.SerialNumber))

	return c.Send(qrCode)
}
