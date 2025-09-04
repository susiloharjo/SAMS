package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"sams-backend/internal/database"
	"sams-backend/internal/handlers"
	"sams-backend/internal/middleware"
	"sams-backend/internal/models"
)

func main() {
	// Removed godotenv.Load() calls to ensure Docker environment variables are used
	// if err := godotenv.Load("../.env"); err != nil {
	// 	log.Println("No root .env file found, trying current directory")
	// 	if err := godotenv.Load(); err != nil {
	// 		log.Println("No .env file found, using system environment variables")
	// 	}
	// }

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate database schema
	if err := db.AutoMigrate(&models.Category{}, &models.Asset{}, &models.Department{}, &models.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize handlers
	userHandler := handlers.NewUserHandler(db)
	authHandler := handlers.NewAuthHandler(db)

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

	// Test login endpoint (temporary for debugging)
	app.Post("/api/v1/test-login", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Test login endpoint working",
			"data": fiber.Map{
				"user": fiber.Map{
					"id":         "test-123",
					"username":   "test",
					"role":       "admin",
					"first_name": "Test",
					"last_name":  "User",
				},
				"access_token":  "test-token-123",
				"refresh_token": "test-refresh-123",
			},
		})
	})

	// Public Authentication Routes (no middleware required)
	app.Post("/api/v1/auth/login", authHandler.Login)
	app.Post("/api/v1/auth/refresh", authHandler.RefreshToken)

	// Protected Routes (require authentication)
	// Asset Routes with RBAC
	app.Get("/api/v1/assets", middleware.AuthMiddleware(), middleware.RequireUser(), handlers.GetAssets)
	app.Get("/api/v1/assets/summary", middleware.AuthMiddleware(), middleware.RequireUser(), handlers.GetAssetSummary)
	app.Get("/api/v1/assets/summary-by-category", middleware.AuthMiddleware(), middleware.RequireUser(), handlers.GetCategorySummary)
	app.Get("/api/v1/assets/summary-by-status", middleware.AuthMiddleware(), middleware.RequireUser(), handlers.GetStatusSummary)
	app.Get("/api/v1/assets/:id", middleware.AuthMiddleware(), middleware.RequireUser(), handlers.GetAsset)
	app.Get("/api/v1/assets/:id/qr", middleware.AuthMiddleware(), middleware.RequireUser(), handlers.GenerateAssetQR)

	// Asset CRUD operations - only admin and manager
	app.Post("/api/v1/assets", middleware.AuthMiddleware(), middleware.RequireManager(), handlers.CreateAsset)
	app.Put("/api/v1/assets/:id", middleware.AuthMiddleware(), middleware.RequireManager(), handlers.UpdateAsset)
	app.Delete("/api/v1/assets/:id", middleware.AuthMiddleware(), middleware.RequireManager(), handlers.DeleteAsset)

	// Category Routes - only admin and manager
	app.Get("/api/v1/categories", middleware.AuthMiddleware(), middleware.RequireUser(), handlers.GetCategories)
	app.Post("/api/v1/categories", middleware.AuthMiddleware(), middleware.RequireManager(), handlers.CreateCategory)
	app.Put("/api/v1/categories/:id", middleware.AuthMiddleware(), middleware.RequireManager(), handlers.UpdateCategory)
	app.Delete("/api/v1/categories/:id", middleware.AuthMiddleware(), middleware.RequireManager(), handlers.DeleteCategory)

	// Department Routes - only admin and manager
	app.Get("/api/v1/departments", middleware.AuthMiddleware(), middleware.RequireUser(), handlers.GetDepartments)
	app.Post("/api/v1/departments", middleware.AuthMiddleware(), middleware.RequireManager(), handlers.CreateDepartment)
	app.Put("/api/v1/departments/:id", middleware.AuthMiddleware(), middleware.RequireManager(), handlers.UpdateDepartment)
	app.Delete("/api/v1/departments/:id", middleware.AuthMiddleware(), middleware.RequireManager(), handlers.DeleteDepartment)

	// AI Routes - only admin and manager
	app.Post("/api/v1/ai/query", middleware.AuthMiddleware(), middleware.RequireManager(), handlers.HandleAIQuery)

	// User Management Routes - only admin
	app.Get("/api/v1/users", middleware.AuthMiddleware(), middleware.RequireAdmin(), userHandler.GetUsers)
	app.Get("/api/v1/users/summary", middleware.AuthMiddleware(), middleware.RequireAdmin(), userHandler.GetUserSummary)
	app.Post("/api/v1/users", middleware.AuthMiddleware(), middleware.RequireAdmin(), userHandler.CreateUser)
	app.Get("/api/v1/users/:id", middleware.AuthMiddleware(), middleware.RequireAdmin(), userHandler.GetUser)
	app.Put("/api/v1/users/:id", middleware.AuthMiddleware(), middleware.RequireAdmin(), userHandler.UpdateUser)
	app.Delete("/api/v1/users/:id", middleware.AuthMiddleware(), middleware.RequireAdmin(), userHandler.DeleteUser)

	// User Profile Routes (authenticated users can manage their own profile)
	app.Post("/api/v1/auth/change-password", middleware.AuthMiddleware(), middleware.RequireUser(), authHandler.ChangePassword)
	app.Post("/api/v1/auth/logout", middleware.AuthMiddleware(), middleware.RequireUser(), authHandler.Logout)

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
