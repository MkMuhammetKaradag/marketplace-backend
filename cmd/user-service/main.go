package main

import (
	application "marketplace/internal/user-service/app"
	"marketplace/internal/user-service/config"
)

func main() {
	appConfig := config.Read()

	app := application.NewApp(appConfig)

	app.Start()
}
