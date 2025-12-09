package middleware

import (
	"fmt"
	"log"
	"marketplace/internal/api-gateway/config"
	"marketplace/internal/api-gateway/limiter"
	"marketplace/internal/api-gateway/metrics"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

func RateLimitMiddleware(rl *limiter.RateLimiter, m *metrics.Metrics, configs map[string]config.RouteConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		clientID := ExtractClientIdentifier(c)
		m.IncrementTotal()
		m.IncrementPath(path)

		routeConfig := getRouteConfig(path, configs)

		// User Limit
		if routeConfig.UserLimit > 0 {
			l := rl.GetLimiter("user:"+clientID+":"+path, toLimit(routeConfig.UserLimit), routeConfig.UserBurst)
			if !l.Allow() {
				m.IncrementRateLimit("user-path")
				log.Printf("⛔ Rate limit (User): %s -> %s", clientID, path)
				c.Set("X-RateLimit-Limit", fmt.Sprintf("%.0f", routeConfig.UserLimit*60))
				c.Set("X-RateLimit-Remaining", fmt.Sprintf("%.0f", l.Tokens()))
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Too many requests",
					"type":  "user-path",
				})
			}
			c.Set("X-RateLimit-Limit", fmt.Sprintf("%.0f", routeConfig.UserLimit*60))
			c.Set("X-RateLimit-Remaining", fmt.Sprintf("%.0f", l.Tokens()))
		}

		// Global Limit
		l := rl.GetLimiter("global:"+path, toLimit(routeConfig.GlobalLimit), routeConfig.GlobalBurst)
		if !l.Allow() {
			m.IncrementRateLimit("global-path")
			log.Printf("⛔ Rate limit (Global): %s", path)
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "System busy",
				"type":  "global-path",
			})
		}

		return c.Next()
	}
}

// Helpers

func toLimit(f float64) rate.Limit {
	return rate.Limit(f)
}

func getRouteConfig(path string, configs map[string]config.RouteConfig) config.RouteConfig {
	conf, exists := configs[path]
	if exists {
		return conf
	}

	// Default fallback
	conf = configs["default"]
	longestMatchLen := 0

	for route, c := range configs {
		if route != "default" && strings.HasPrefix(path, route) {
			if len(route) > longestMatchLen {
				conf = c
				longestMatchLen = len(route)
			}
		}
	}
	return conf
}

func ExtractClientIdentifier(c *fiber.Ctx) string {
	if cookie := c.Cookies(config.SessionCookieName); cookie != "" {
		return "session:" + cookie
	}
	authHeader := c.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") && len(authHeader) > 7 {
		return "token:" + authHeader[7:15]
	}
	ip := c.IP()
	return "ip:" + ip
}
