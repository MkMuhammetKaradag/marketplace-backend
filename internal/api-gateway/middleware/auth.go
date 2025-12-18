package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"marketplace/internal/api-gateway/config"

	"marketplace/internal/api-gateway/grpc_client"
)

// AuthMiddleware checks for session cookie or authorization header
func AuthMiddleware(policies map[string]config.RoutePolicy) fiber.Handler {
	return func(c *fiber.Ctx) error {

		routePath := c.Route().Path
		requestPath := c.Path()

		policy, isProtected := policies[routePath]

		if !isProtected {
			isProtected, policy = findParametrizedRoute(requestPath, policies)
		}

		if !isProtected {
			return c.Next()
		}

		var authValue string
		if cookie := c.Cookies(config.SessionCookieName); cookie != "" {
			authValue = cookie
		}
		if authHeader := c.Get("Authorization"); strings.HasPrefix(authHeader, "Bearer ") {
			authValue = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if authValue == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		var userID string
		var role string
		// var userRole string
		var isValid bool

		isValid, userID, role = grpc_client.ValidateToken(authValue)

		if !isValid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid session/token",
			})
		}

		if !contains(policy.Roles, role) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You do not have permission to access this resource",
			})
		}

		c.Locals("userID", userID)
		c.Locals("role", role)
		c.Request().Header.Set("X-User-ID", userID)
		c.Request().Header.Set("X-User-Role", role)

		return c.Next()
	}
}

// findParametrizedRoute, isteği haritadaki parametreli şablonlarla eşleştirmeye çalışır
func findParametrizedRoute(requestPath string, policies map[string]config.RoutePolicy) (bool, config.RoutePolicy) {
	requestParts := strings.Split(requestPath, "/")

	for policyPath, policy := range policies {
		policyParts := strings.Split(policyPath, "/")

		if len(policyParts) != len(requestParts) {
			continue
		}

		isMatch := true
		for i := 0; i < len(policyParts); i++ {
			policyPart := policyParts[i]
			requestPart := requestParts[i]

			if strings.HasPrefix(policyPart, ":") || policyPart == "" {
				continue
			}

			if policyPart != requestPart {
				isMatch = false
				break
			}
		}

		if isMatch {
			return true, policy
		}
	}
	return false, config.RoutePolicy{}
}

func contains(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
