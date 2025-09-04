package middleware

import (
	"errors"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request is coming from MCP server (internal container)
		if c.IP() == "172.18.0.5" || c.IP() == "172.18.0.6" {
			// Allow MCP server requests without authentication
			c.Locals("user_id", "cd7874a9-a08d-4292-8bef-e94d64c1ceb7") // Admin user ID
			c.Locals("role", "admin")
			return c.Next()
		}

		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Authorization header required",
			})
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid authorization header format",
			})
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			// Get JWT secret from environment
			jwtSecret := os.Getenv("JWT_SECRET")
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid or expired token",
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

		// Set user information in context
		userID, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid user ID in token",
			})
		}

		username, ok := claims["username"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid username in token",
			})
		}

		role, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid role in token",
			})
		}

		// Set user context
		c.Locals("user_id", userID)
		c.Locals("username", username)
		c.Locals("role", role)

		return c.Next()
	}
}

// RequireRole creates middleware that requires specific roles
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user role from context (set by AuthMiddleware)
		userRole := c.Locals("role").(string)
		if userRole == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "User not authenticated",
			})
		}

		// Check if user has required role
		hasRole := false
		for _, role := range allowedRoles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// RequireAdmin creates middleware that requires admin role
func RequireAdmin() fiber.Handler {
	return RequireRole("admin")
}

// RequireManager creates middleware that requires manager or admin role
func RequireManager() fiber.Handler {
	return RequireRole("admin", "manager")
}

// RequireUser creates middleware that requires any authenticated user
func RequireUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user role from context (set by AuthMiddleware)
		userRole := c.Locals("role").(string)
		if userRole == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "User not authenticated",
			})
		}

		return c.Next()
	}
}

// GetCurrentUserID returns the current user ID from context
func GetCurrentUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userIDStr := c.Locals("user_id").(string)
	if userIDStr == "" {
		return uuid.Nil, errors.New("user not authenticated")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

// GetCurrentUserRole returns the current user role from context
func GetCurrentUserRole(c *fiber.Ctx) string {
	role := c.Locals("role").(string)
	return role
}

// CanPerformAction checks if user can perform specific actions based on role
func CanPerformAction(c *fiber.Ctx, action string, resourceOwnerID *uuid.UUID) bool {
	userRole := GetCurrentUserRole(c)
	_, err := GetCurrentUserID(c)
	if err != nil {
		return false
	}

	// Admin can do everything
	if userRole == "admin" {
		return true
	}

	// Manager can do most things
	if userRole == "manager" {
		return true
	}

	// Regular users have limited permissions
	if userRole == "user" {
		switch action {
		case "read":
			// Users can read all assets
			return true
		case "create":
			// Users cannot create assets
			return false
		case "update":
			// Users cannot update assets
			return false
		case "delete":
			// Users cannot delete assets
			return false
		default:
			return false
		}
	}

	return false
}

// AssetAccessControl checks if user can access specific asset operations
func AssetAccessControl(action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := GetCurrentUserRole(c)

		// Admin and manager have full access
		if userRole == "admin" || userRole == "manager" {
			return c.Next()
		}

		// Regular users can only read
		if userRole == "user" && action == "read" {
			return c.Next()
		}

		// For other actions, users are denied
		if userRole == "user" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "Insufficient permissions for this operation",
			})
		}

		return c.Next()
	}
}
