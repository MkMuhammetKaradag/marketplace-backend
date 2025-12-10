package middleware

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"

	"marketplace/internal/api-gateway/config"

	"marketplace/internal/api-gateway/grpc_client"
)

// AuthMiddleware checks for session cookie or authorization header
func AuthMiddleware(protectedPaths map[string]bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		// Check if path is protected
		if protectedPaths[path] {
			var authValue string // Cookie veya Token değeri

			// 1. Cookie Check
			cookie := c.Cookies(config.SessionCookieName)
			if cookie != "" {
				authValue = cookie
			}

			// 2. Authorization Header Check (Bearer Token)
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") && len(authHeader) > 7 {
				authValue = strings.TrimPrefix(authHeader, "Bearer ")
			}

			// Eğer bir token/cookie bulunduysa, gRPC ile User Servisine sor
			isAuthenticated := false
			var userID string

			if authValue != "" {
				isValid, id := grpc_client.ValidateToken(authValue)
				fmt.Println("isValid", isValid)
				fmt.Println("id", id)

				if isValid {
					isAuthenticated = true
					userID = id
				}
			}

			if !isAuthenticated {
				log.Printf("🔒 Unauthorized access: %s", path)
				// Token/Cookie vardı ama User Servisi geçersiz dedi.
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Authentication (Session/Token) required or invalid",
				})
			}

			c.Request().Header.Set("X-User-ID", userID)
		}

		return c.Next()
	}
}
