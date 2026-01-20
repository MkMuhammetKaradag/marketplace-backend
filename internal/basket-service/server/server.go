// internal/basket-service/server/server.go
package server

import (
	"fmt"
	"log"
	"marketplace/internal/basket-service/domain"
	"net"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"google.golang.org/grpc"
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
	app               *fiber.App
	cfg               Config
	grpcProductClient domain.ProductClient
	grpcServer        *grpc.Server
}
type GrpcServerRegistrar interface {
	Register(server *grpc.Server)
}

func New(cfg Config, registrar RouteRegistrar, pClient domain.ProductClient, grpcHandler GrpcServerRegistrar) *Server {
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
	grpcSrv := grpc.NewServer()

	if grpcHandler != nil {
		grpcHandler.Register(grpcSrv)
	}
	return &Server{
		app:               app,
		cfg:               cfg,
		grpcProductClient: pClient,
		grpcServer:        grpcSrv,
	}
}

func (s *Server) startGrpc() error {
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", s.cfg.GrpcPort))
	if err != nil {
		return fmt.Errorf("gRPC dinlemede hata: %w", err)
	}
	log.Printf("👂 gRPC sunucusu %s adresinde dinliyor...", s.cfg.GrpcPort)

	return s.grpcServer.Serve(listen)
}

func (s *Server) Start() error {
	go func() {
		if err := s.startGrpc(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("gRPC sunucusu hatası: %v", err)
		}
	}()
	log.Printf("🌐 HTTP sunucusu %s adresinde dinliyor...", s.cfg.Port)
	return s.app.Listen(s.Address())
}

func (s *Server) Shutdown(timeout time.Duration) error {
	if s.grpcProductClient != nil {
		s.grpcProductClient.Close()
	}
	s.grpcServer.GracefulStop()
	return s.app.ShutdownWithTimeout(timeout)
}

func (s *Server) FiberApp() *fiber.App {
	return s.app
}

func (s *Server) Address() string {
	return fmt.Sprintf("0.0.0.0:%s", s.cfg.Port)
}

// func (s *Server) Run() error {
// 	grpcAddress := "localhost:3004" // Docker'da ise servis adı, yerelde ise localhost:50051

// 	if err := s.grpcProductClient.InitProductServiceClient(grpcAddress); err != nil {
// 		log.Fatalf("gRPC istemcisi başlatılamadı: %v", err)
// 		return err
// 	}
// 	return nil
// }
