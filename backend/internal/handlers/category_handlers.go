package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"sams-backend/internal/database"
	"sams-backend/internal/models"
)

func GetCategories(c *fiber.Ctx) error {
	db := database.GetDB()
	var categories []models.Category
	if err := db.Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to fetch categories"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "data": categories})
}

func GetCategory(c *fiber.Ctx) error {
	db := database.GetDB()
	var category models.Category
	id := c.Params("id")

	if err := db.First(&category, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "message": "Category not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to fetch category"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "data": category})
}

func CreateCategory(c *fiber.Ctx) error {
	db := database.GetDB()
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid request body"})
	}

	if err := db.Create(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to create category"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"error": false, "data": category})
}

func UpdateCategory(c *fiber.Ctx) error {
	db := database.GetDB()
	var category models.Category
	id := c.Params("id")

	if err := db.First(&category, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "message": "Category not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to fetch category"})
	}

	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid request body"})
	}

	if err := db.Save(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to update category"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "data": category})
}

func DeleteCategory(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")

	result := db.Delete(&models.Category{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to delete category"})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "message": "Category not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "message": "Category deleted successfully"})
}
