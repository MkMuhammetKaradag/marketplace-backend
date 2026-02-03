// internal/notification-service/server/server.go
package server

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type RouteRegistrar interface {
	Register(app *fiber.App)
}

type Config struct {
	GrpcPort     string
	Port         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Server struct {
	app *fiber.App
	cfg Config
}

func New(cfg Config, registrar RouteRegistrar) *Server {
	app := fiber.New(fiber.Config{
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Concurrency:  256 * 1024,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowCredentials: true,
	}))
	app.Use(requestid.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "UP"})
	})

	if registrar != nil {
		registrar.Register(app)
	}

	return &Server{
		app: app,
		cfg: cfg,
	}
}

func (s *Server) Start() error {
	log.Printf("🌐 HTTP sunucusu %s adresinde dinliyor...", s.cfg.Port)
	return s.app.Listen(s.Address())
}

func (s *Server) Shutdown(timeout time.Duration) error {
	return s.app.ShutdownWithTimeout(timeout)
}

func (s *Server) FiberApp() *fiber.App {
	return s.app
}

func (s *Server) Address() string {
	return fmt.Sprintf("0.0.0.0:%s", s.cfg.Port)
}
