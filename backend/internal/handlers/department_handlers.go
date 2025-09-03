package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"sams-backend/internal/database"
	"sams-backend/internal/models"
)

func GetDepartments(c *fiber.Ctx) error {
	db := database.GetDB()
	var departments []models.Department
	if err := db.Find(&departments).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to fetch departments"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "data": departments})
}

func GetDepartment(c *fiber.Ctx) error {
	db := database.GetDB()
	var department models.Department
	id := c.Params("id")

	if err := db.First(&department, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "message": "Department not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to fetch department"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "data": department})
}

func CreateDepartment(c *fiber.Ctx) error {
	db := database.GetDB()
	var department models.Department
	if err := c.BodyParser(&department); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid request body"})
	}

	if err := db.Create(&department).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to create department"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"error": false, "data": department})
}

func UpdateDepartment(c *fiber.Ctx) error {
	db := database.GetDB()
	var department models.Department
	id := c.Params("id")

	if err := db.First(&department, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "message": "Department not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to fetch department"})
	}

	if err := c.BodyParser(&department); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid request body"})
	}

	if err := db.Save(&department).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to update department"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "data": department})
}

func DeleteDepartment(c *fiber.Ctx) error {
	db := database.GetDB()
	id := c.Params("id")

	result := db.Delete(&models.Department{}, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to delete department"})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "message": "Department not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "message": "Department deleted successfully"})
}
