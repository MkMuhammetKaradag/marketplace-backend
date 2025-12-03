// internal/user-service/server/server.go
package app

import (
	"fmt"
	"marketplace/internal/user-service/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Config struct {
	Port         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func SetupServer(config config.Config, httpHandlers map[string]interface{}) *fiber.App {

	serverConfig := Config{
		Port:         config.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	myApp := NewFiberApp(serverConfig)

	myApp.Get("/hello", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"message": "Hello from User Service!",
			"info":    "Header kontrol et",
		})
	})
	return myApp
}

func NewFiberApp(cfg Config) *fiber.App {
	myApp := fiber.New(fiber.Config{
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Concurrency:  256 * 1024,
	})
	myApp.Use(cors.New(cors.Config{
		// Frontend'inin tam adresini buraya yaz.
		// Bu, * yerine geçen ve kimlik bilgilerini göndermeye izin veren adrestir.
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowCredentials: true,
	}))
	myApp.Use(requestid.New())

	// basic health endpoint
	myApp.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "UP"})
	})
	return myApp
}

func Start(myApp *fiber.App, port string) error {
	return myApp.Listen(fmt.Sprintf("0.0.0.0:%s", port))
}
