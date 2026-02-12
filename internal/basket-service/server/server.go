package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"google.golang.org/grpc"
)

type RouteRegistrar interface {
	Register(app *fiber.App)
}

type GrpcServerRegistrar interface {
	Register(server *grpc.Server)
}

type Config struct {
	GrpcPort     string
	Port         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Server struct {
	app        *fiber.App
	cfg        Config
	grpcServer *grpc.Server
}

func New(cfg Config, registrar RouteRegistrar, grpcHandler GrpcServerRegistrar) *Server {
	app := fiber.New(fiber.Config{
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	app.Use(requestid.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowCredentials: true,
	}))

	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "UP", "service": "basket-service"})
	})

	// HTTP Rotalarını Kaydet
	if registrar != nil {
		registrar.Register(app)
	}

	// gRPC Sunucusu Kurulumu
	grpcSrv := grpc.NewServer()
	if grpcHandler != nil {
		grpcHandler.Register(grpcSrv)
	}

	return &Server{
		app:        app,
		cfg:        cfg,
		grpcServer: grpcSrv,
	}
}

func (s *Server) Start() error {
	// gRPC'yi ayrı bir goroutine'de başlat
	go func() {
		if err := s.startGrpc(); err != nil {
			log.Printf("❌ gRPC sunucusu hatası: %v", err)
		}
	}()

	log.Printf("🌐 HTTP sunucusu %s adresinde dinliyor...", s.cfg.Port)
	return s.app.Listen(s.Address())
}

func (s *Server) startGrpc() error {
	addr := fmt.Sprintf("0.0.0.0:%s", s.cfg.GrpcPort)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("gRPC listener failed: %w", err)
	}
	log.Printf("👂 gRPC sunucusu %s adresinde dinliyor...", addr)
	return s.grpcServer.Serve(listen)
}

func (s *Server) Shutdown(timeout time.Duration) error {
	log.Println("Shutting down servers...")

	// gRPC'yi nazikçe durdur
	s.grpcServer.GracefulStop()

	// Fiber'i durdur
	return s.app.ShutdownWithTimeout(timeout)
}

func (s *Server) FiberApp() *fiber.App {
	return s.app
}

func (s *Server) Address() string {
	return fmt.Sprintf("0.0.0.0:%s", s.cfg.Port)
}
