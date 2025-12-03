// internal/user-service/app/application.go
package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/user-service/config"
	"marketplace/internal/user-service/database/postgres"
	"marketplace/internal/user-service/domain"
	graceful "marketplace/internal/user-service/pkg/gareceful"
	"time"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	config       config.Config
	postgresRepo domain.PostgresRepository
	fiberApp     *fiber.App
	httpHandlers map[string]interface{}
}

func NewApp(config config.Config) *App {
	myApp := &App{
		config: config,
	}
	myApp.initDependencies()
	return myApp
}

func (a *App) initDependencies() {

	repo, err := postgres.NewRepository(a.config)
	if err != nil {
		log.Fatal("DB init failed:", err)
	}
	a.postgresRepo = repo
	a.httpHandlers = SetupHTTPHandlers()
	a.fiberApp = SetupServer(a.config, a.httpHandlers)
}

func (a *App) Start() {
	go func() {
		port := a.config.Server.Port
		if err := a.fiberApp.Listen(":" + port); err != nil {
			fmt.Println("Failed to start server:", err)
		}
	}()

	graceful.WaitForShutdown(a.fiberApp, 5*time.Second, context.Background())
}
