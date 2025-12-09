package handlers

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/valyala/fasthttp"

	"marketplace/internal/api-gateway/config"
	"marketplace/internal/api-gateway/metrics"
	"marketplace/internal/api-gateway/service"
)

type ProxyHandler struct {
	Registry *service.ServiceRegistry
	Metrics  *metrics.Metrics
}

func NewProxyHandler(registry *service.ServiceRegistry, metrics *metrics.Metrics) *ProxyHandler {
	return &ProxyHandler{
		Registry: registry,
		Metrics:  metrics,
	}
}

func (h *ProxyHandler) Handle(c *fiber.Ctx) error {
	path := c.Path()
	svc, ok := h.Registry.GetByPath(path)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Service not found"})
	}

	// Circuit Breaker Check
	if !h.Registry.IsHealthy(svc) {
		h.Metrics.IncrementFailed()
		log.Printf("❌ Circuit Open: %s", svc.Name)
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":   "Service unavailable (Circuit Open)",
			"service": svc.Name,
		})
	}

	targetBaseURL, ok := svc.GetNextBaseURL()
	if !ok {
		h.Metrics.IncrementFailed()
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":   "No healthy service instance found",
			"service": svc.Name,
		})
	}

	h.Metrics.IncrementService(svc.Name)

	targetPath := strings.TrimPrefix(path, svc.PathPrefix)
	targetURL := targetBaseURL + targetPath
	if len(c.Request().URI().QueryString()) > 0 {
		targetURL += "?" + string(c.Request().URI().QueryString())
	}

	// Setup Request Headers
	c.Request().Header.Set(config.InternalGatewayHeader, config.InternalGatewaySecret)
	c.Request().Header.Set("X-Forwarded-For", c.IP())

	// Fiber Proxy Middleware `Do`
	// Note: proxy.Do uses fasthttp.Client for the client argument
	if err := proxy.Do(c, targetURL, &fasthttp.Client{
		ReadTimeout:  svc.Timeout,
		WriteTimeout: svc.Timeout,
	}); err != nil {
		h.Metrics.IncrementFailed()
		log.Printf("❌ Proxy error [%s]: %v", svc.Name, err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Backend service error"})
	}

	h.Metrics.IncrementSuccess()
	return nil
}
