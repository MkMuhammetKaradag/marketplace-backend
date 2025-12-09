package middleware

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"

	"marketplace/internal/api-gateway/config"
)

// AuthMiddleware checks for session cookie or authorization header
func AuthMiddleware(protectedPaths map[string]bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		// Check if path is protected
		if protectedPaths[path] {
			isAuthenticated := false

			// 1. Cookie Check
			cookie := c.Cookies(config.SessionCookieName)
			if cookie != "" {
				isAuthenticated = true
			}

			// 2. Authorization Header Check (Bearer Token)
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") && len(authHeader) > 7 {
				isAuthenticated = true
			}

			if !isAuthenticated {
				log.Printf("🔒 Unauthorized access: %s", path)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Authentication (Session/Token) required",
				})
			}
		}
		return c.Next()
	}
}
