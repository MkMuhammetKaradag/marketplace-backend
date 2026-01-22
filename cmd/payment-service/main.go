package main

import (
	"log"

	application "marketplace/internal/payment-service/app"
	"marketplace/internal/payment-service/config"
)

func main() {
	appConfig := config.Read()

	app, err := application.NewApp(appConfig)
	if err != nil {
		log.Fatalf("failed to initialise app: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
