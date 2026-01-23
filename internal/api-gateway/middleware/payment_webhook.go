package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func WebhookSecurityMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		if strings.HasSuffix(path, "/payment/webhook") {
			sigHeader := c.Get("Stripe-Signature")
			fmt.Println(sigHeader)

			if sigHeader == "" {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Direct access to webhook is not allowed",
				})
			}

		}

		return c.Next()
	}
}
