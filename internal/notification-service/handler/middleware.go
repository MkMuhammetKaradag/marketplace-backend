// internal/notification-service/handler/middleware.go
package handler

import (
	"errors"
	"marketplace/internal/notification-service/domain"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func HandleBasic[R Request, Res Response](handler BasicHandler[R, Res]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req R

		if err := parseRequest(c, &req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed", "details": err.Error()})
		}

		ctx := c.UserContext()
		res, err := handler.Handle(ctx, &req)

		if err != nil {
			status := getStatusCodeFromError(err)
			return c.Status(status).JSON(fiber.Map{"error": err.Error()})
		}
		if c.Method() == fiber.MethodPost {
			return c.Status(fiber.StatusCreated).JSON(res)
		}
		return c.JSON(res)
	}
}
func HandleWithFiber[R Request, Res Response](handler FiberHandler[R, Res]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req R

		if err := parseRequest(c, &req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed", "details": err.Error()})
		}

		// ctx := c.UserContext()
		res, err := handler.Handle(c, &req)

		if err != nil {
			status := getStatusCodeFromError(err)
			return c.Status(status).JSON(fiber.Map{"error": err.Error()})
		}
		if c.Method() == fiber.MethodPost {
			return c.Status(fiber.StatusCreated).JSON(res)
		}
		return c.JSON(res)
	}
}
func parseRequest[R any](c *fiber.Ctx, req *R) error {
	if err := c.BodyParser(req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
		return err
	}

	if err := c.ParamsParser(req); err != nil {
		return err
	}

	if err := c.QueryParser(req); err != nil {
		return err
	}

	if err := c.ReqHeaderParser(req); err != nil {
		return err
	}

	return nil
}
func getStatusCodeFromError(err error) int {
	switch {
	case errors.Is(err, domain.ErrOrderNotFound):
		return fiber.StatusNotFound
	case errors.Is(err, domain.ErrUnauthorized):
		return fiber.StatusUnauthorized
	default:
		return fiber.StatusInternalServerError
	}
}
