package main

import (
	"log"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"sams-backend/internal/database"
	"sams-backend/internal/models"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate database schema
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create default users
	users := []models.User{
		{
			Username:  "admin",
			Email:     "admin@sams.com",
			FirstName: "System",
			LastName:  "Administrator",
			Password:  "admin123",
			Role:      "admin",
			IsActive:  true,
		},
		{
			Username:  "manager",
			Email:     "manager@sams.com",
			FirstName: "System",
			LastName:  "Manager",
			Password:  "manager123",
			Role:      "manager",
			IsActive:  true,
		},
		{
			Username:  "user",
			Email:     "user@sams.com",
			FirstName: "System",
			LastName:  "User",
			Password:  "user123",
			Role:      "user",
			IsActive:  true,
		},
	}

	for _, user := range users {
		// Check if user already exists
		var existingUser models.User
		if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
			log.Printf("User %s already exists, skipping...", user.Username)
			continue
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for user %s: %v", user.Username, err)
			continue
		}

		// Set hashed password
		user.Password = string(hashedPassword)

		// Create user
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", user.Username, err)
			continue
		}

		log.Printf("Successfully created user: %s (%s)", user.Username, user.Role)
	}

	log.Println("User creation completed!")
	log.Println("\nDefault credentials:")
	log.Println("Admin: username=admin, password=admin123")
	log.Println("Manager: username=manager, password=manager123")
	log.Println("User: username=user, password=user123")
}
