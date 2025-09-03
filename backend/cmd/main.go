package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"sams-backend/internal/database"
	"sams-backend/internal/handlers"
	"sams-backend/internal/models"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate database schema
	if err := db.AutoMigrate(&models.Category{}, &models.Asset{}, &models.Department{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "SAMS Backend",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())

	// Basic CORS headers for MVP
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	})

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "SAMS Backend",
			"version": "1.0.0",
		})
	})

	// API Routes
	app.Get("/api/v1/categories", handlers.GetCategories)
	app.Post("/api/v1/categories", handlers.CreateCategory)
	app.Put("/api/v1/categories/:id", handlers.UpdateCategory)
	app.Delete("/api/v1/categories/:id", handlers.DeleteCategory)

	app.Get("/api/v1/departments", handlers.GetDepartments)
	app.Post("/api/v1/departments", handlers.CreateDepartment)
	app.Put("/api/v1/departments/:id", handlers.UpdateDepartment)
	app.Delete("/api/v1/departments/:id", handlers.DeleteDepartment)

	app.Get("/api/v1/assets", handlers.GetAssets)
	app.Get("/api/v1/assets/summary", handlers.GetAssetSummary) // Moved up
	app.Get("/api/v1/assets/summary-by-category", handlers.GetCategorySummary)
	app.Post("/api/v1/assets", handlers.CreateAsset)
	app.Get("/api/v1/assets/:id", handlers.GetAsset)
	app.Put("/api/v1/assets/:id", handlers.UpdateAsset)
	app.Delete("/api/v1/assets/:id", handlers.DeleteAsset)
	app.Get("/api/v1/assets/:id/qr", handlers.GenerateAssetQR)

	app.Post("/api/v1/ai/query", handlers.HandleAIQuery)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting SAMS Backend on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
