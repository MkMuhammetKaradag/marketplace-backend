package middleware

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"marketplace/internal/api-gateway/cache"
	"marketplace/internal/api-gateway/config"

	"marketplace/internal/api-gateway/grpc_client"
)

const (
	PermissionAdministrator int64 = 1 << 62
)

// AuthMiddleware checks for session cookie or authorization header
func AuthMiddleware(policies map[string]config.RoutePolicy, cacheManager *cache.CacheManager) fiber.Handler {
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
		var permissions int64
		// var userRole string
		var isValid bool
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		cachedSession, cacheErr := cacheManager.GetSession(ctx, authValue)

		if cacheErr == nil && cachedSession != nil {
			// Cache'de bulundu! ✅
			log.Printf("✅ Cache HIT - UserID: %s", cachedSession.UserID)
			userID = cachedSession.UserID
			permissions = cachedSession.Permissions
			isValid = true
		} else {
			isValid, userID, permissions = grpc_client.ValidateToken(authValue)

			if !isValid {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid session/token",
				})
			}
			if err := cacheManager.SetSession(ctx, authValue, userID, permissions); err != nil {
				log.Printf("⚠️ Cache save error: %v", err)
				// Hata olsa bile kullanıcıyı engelleme
			} else {
				log.Printf("💾 Cache save success - UserID: %s", userID)
			}
		}
		if !contains(policy.Permissions, permissions) {
			log.Printf("❌ Permission denied - UserID: %s", userID)
			log.Printf("❌ Permission denied - Permissions: %s", permissions)
			log.Printf("❌ Permission denied - Required Permissions: %s", policy.Permissions)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You do not have permission to access this resource",
			})
		}

		c.Locals("userID", userID)
		c.Locals("permissions", permissions)
		c.Request().Header.Set("X-User-ID", userID)
		c.Request().Header.Set("X-User-Permissions", strconv.FormatInt(permissions, 10))

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

func contains(requiredPerm int64, userTotalPerms int64) bool {
	if (userTotalPerms & PermissionAdministrator) == PermissionAdministrator {
		return true
	}
	return (userTotalPerms & requiredPerm) == requiredPerm
}
