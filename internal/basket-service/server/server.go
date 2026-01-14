// internal/basket-service/server/server.go
package server

import (
	"fmt"
	"log"
	"marketplace/internal/basket-service/grpc_client"
	"net/http"
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
	go func() {
		if err := s.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("gRPC sunucusu hatası: %v", err)
		}
	}()
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
func (s *Server) Run() error {
	grpcAddress := "localhost:3004" // Docker'da ise servis adı, yerelde ise localhost:50051

	if err := grpc_client.InitProductServiceClient(grpcAddress); err != nil {
		log.Fatalf("gRPC istemcisi başlatılamadı: %v", err)
		return err
	}
	return nil
}
