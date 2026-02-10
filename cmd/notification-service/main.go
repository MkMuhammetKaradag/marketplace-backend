package main

import (
    "log"
    "marketplace/internal/notification-service/app" 
    "marketplace/internal/notification-service/config"
)

func main() {
    // 1. Konfigürasyonu yükle
    appConfig := config.Read()

    // 2. Uygulamayı ayağa kaldır (NewApp kullanımı daha yaygındır)
    application, err := app.NewApp(appConfig) 
    if err != nil {
        log.Fatalf("failed to initialise app: %v", err)
    }

    // 3. Uygulamayı başlat
    if err := application.Start(); err != nil {
        log.Fatalf("server stopped with error: %v", err)
    }
}