package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"marketplace/internal/api-gateway/config"
	"marketplace/internal/api-gateway/metrics"
	"marketplace/internal/api-gateway/service"
)

type ManageHandler struct {
	Registry *service.ServiceRegistry
	Metrics  *metrics.Metrics
}

func NewManageHandler(registry *service.ServiceRegistry, metrics *metrics.Metrics) *ManageHandler {
	return &ManageHandler{
		Registry: registry,
		Metrics:  metrics,
	}
}

func (h *ManageHandler) HealthCheck(c *fiber.Ctx) error {
	services := h.Registry.List()
	serviceHealth := make(map[string]interface{})

	for _, svc := range services {
		healthy := h.Registry.IsHealthy(svc)
		// We need to access FailCount and LastCheck but they are not easily accessible thread-safely
		// if we strictly follow encapsulation.
		// For now, let's just return basic info or access directly if we are ok with race (for reading stats it is usually ok).
		// Or better, let's stick to what we can access.
		serviceHealth[svc.Name] = map[string]interface{}{
			"healthy":   healthy,
			"base_urls": svc.BaseURLs,
		}
	}

	return c.JSON(fiber.Map{
		"gateway":   "healthy",
		"services":  serviceHealth,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (h *ManageHandler) GetMetrics(c *fiber.Ctx) error {
	stats := h.Metrics.GetStats()
	return c.JSON(stats)
}

func (h *ManageHandler) ListServices(c *fiber.Ctx) error {
	services := h.Registry.List()
	serviceList := make([]map[string]interface{}, 0, len(services))

	for _, svc := range services {
		serviceList = append(serviceList, map[string]interface{}{
			"name":        svc.Name,
			"base_urls":   svc.BaseURLs,
			"path_prefix": svc.PathPrefix,
			"healthy":     h.Registry.IsHealthy(svc),
		})
	}

	return c.JSON(serviceList)
}

func (h *ManageHandler) SimulateLogin(c *fiber.Ctx) error {
	sessionID := uuid.New().String()

	c.Cookie(&fiber.Cookie{
		Name:     config.SessionCookieName,
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false, // Set true in production
	})

	return c.JSON(fiber.Map{
		"message":    "Session created successfully",
		"session_id": sessionID,
		"warning":    "This is a simulation.",
	})
}
