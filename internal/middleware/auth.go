package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"awesomeEcommerce/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware that checks if the request has a valid API key
func AuthMiddleware(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Extract the token
		token := parts[1]

		// In a real application, we would validate the token against a database or cache
		// For now, we'll just check if it's a valid format (simple hash)
		if len(token) < 10 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// For demonstration purposes, we'll extract user information from the token
		// In a real application, we would decode the token and validate it
		// Here we're just simulating by using the token as an email lookup
		user, err := userService.GetUserByEmail(c, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set the user ID and role in the context
		c.Set("userID", user.ID)
		c.Set("email", user.Email)
		c.Set("role", user.Role)

		c.Next()
	}
}

// RoleMiddleware is a middleware that checks if the user has the required role
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user role from the context
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Check if the user has one of the required roles
		userRole := role.(string)
		for _, r := range roles {
			if r == userRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		c.Abort()
	}
}

// GenerateToken generates a simple token for a user
// In a real application, we would use JWT or another token format
func GenerateToken(email string) string {
	// Create a simple token based on email and timestamp
	data := fmt.Sprintf("%s:%d", email, time.Now().Unix())
	hasher := sha256.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

// LoginHandler handles user login and returns a token
func LoginHandler(c *gin.Context, userService service.UserService) {
	// Parse request
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate user
	user, err := userService.AuthenticateUser(c, request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token
	token := GenerateToken(user.Email)

	// Return token
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
