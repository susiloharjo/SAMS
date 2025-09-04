package handlers

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"sams-backend/internal/models"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Login handles user authentication and returns JWT tokens
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Find user by username
	var user models.User
	if err := h.db.Preload("Department").Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid credentials",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to authenticate user",
		})
	}

	// Check if user is active
	if !user.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Account is deactivated",
		})
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid credentials",
		})
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now

	// Use Update instead of Save to only update the LastLogin field
	// This prevents foreign key constraint issues with department_id
	h.db.Model(&user).UpdateColumn("last_login", &now)

	// Generate JWT token
	accessToken, err := h.generateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to generate access token",
		})
	}

	// Generate refresh token
	refreshToken, err := h.generateRefreshToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to generate refresh token",
		})
	}

	// Get token expiration
	expirationTime := time.Now().Add(24 * time.Hour) // 24 hours

	// Convert user to response format
	userResponse := models.UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role,
		DepartmentID: user.DepartmentID,
		Department:   user.Department,
		IsActive:     user.IsActive,
		LastLogin:    user.LastLogin,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	response := models.LoginResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expirationTime,
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"data":    response,
		"message": "Login successful",
	})
}

// RefreshToken generates a new access token using a refresh token
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Parse and validate refresh token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid refresh token",
		})
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid token claims",
		})
	}

	// Get user ID from claims
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid token claims",
		})
	}

	// Find user
	var user models.User
	if err := h.db.First(&user, "id = ?", userIDStr).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "User not found",
		})
	}

	// Check if user is still active
	if !user.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Account is deactivated",
		})
	}

	// Generate new access token
	accessToken, err := h.generateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to generate access token",
		})
	}

	// Generate new refresh token
	refreshToken, err := h.generateRefreshToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to generate refresh token",
		})
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	response := models.LoginResponse{
		User:         models.UserResponse{}, // Don't send full user data on refresh
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expirationTime,
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"data":    response,
		"message": "Token refreshed successfully",
	})
}

// ChangePassword allows users to change their password
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	userID := c.Locals("user_id").(string)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "User not authenticated",
		})
	}

	var req models.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Find user
	var user models.User
	if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "User not found",
		})
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Current password is incorrect",
		})
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to hash password",
		})
	}

	// Update password
	user.Password = string(hashedPassword)
	if err := h.db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update password",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Password changed successfully",
	})
}

// Logout handles user logout (client-side token removal)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// In a stateless JWT system, logout is handled client-side
	// The server can't invalidate JWT tokens, but we can return success
	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Logout successful",
	})
}

// generateJWT creates a new JWT access token
func (h *AuthHandler) generateJWT(user models.User) (string, error) {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")

	// Get token expiration from environment (default: 24 hours)
	expirationStr := os.Getenv("JWT_EXPIRATION_HOURS")
	expirationHours := 24
	if expirationStr != "" {
		if parsed, err := strconv.Atoi(expirationStr); err == nil {
			expirationHours = parsed
		}
	}

	// Create claims using jwt.MapClaims for compatibility
	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)
	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"role":     user.Role,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// generateRefreshToken creates a new JWT refresh token
func (h *AuthHandler) generateRefreshToken(user models.User) (string, error) {
	// Get JWT refresh secret from environment
	jwtRefreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if jwtRefreshSecret == "" {
		// If no specific refresh secret is set, use the same as the main JWT secret
		jwtRefreshSecret = os.Getenv("JWT_SECRET")
	}

	// Refresh tokens last longer (default: 7 days)
	expirationStr := os.Getenv("JWT_REFRESH_EXPIRATION_DAYS")
	expirationDays := 7
	if expirationStr != "" {
		if parsed, err := strconv.Atoi(expirationStr); err == nil {
			expirationDays = parsed
		}
	}

	// Create claims
	expirationTime := time.Now().AddDate(0, 0, expirationDays)
	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"role":     user.Role,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
		"type":     "refresh",
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtRefreshSecret))
}
