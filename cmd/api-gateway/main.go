package main

import (
	"log"

	"marketplace/internal/api-gateway/app"
	"marketplace/internal/api-gateway/config"
)

func main() {
	// Initialize Application
	application := app.New()

	// Register Services (In a real scenario, this might be dynamic or loaded from config)
	// Important: Ensure these URLs match your actual service ports or docker-compose
	application.RegisterService("user-service", []string{"http://localhost:8081", "http://localhost:8081"}, "/users")
	application.RegisterService("seller-service", []string{"http://localhost:8083"}, "/sellers")
	application.RegisterService("auth-service", []string{"http://localhost:8084"}, "/auth")
	application.RegisterService("test-service", []string{"http://localhost:8082"}, "/test")
	application.RegisterService("chat-service", []string{"http://localhost:8085"}, "/chat")

	log.Printf("🚀 Gateway started on %s", config.GatewayPort)
	log.Printf("ℹ️  Usage:")
	log.Printf("  - /users/profile -> user-service (Auth required)")
	log.Printf("  - /test/hello    -> test-service (Strict Rate Limit)")
	log.Printf("  - /simulate/login -> Create test session")

	// Start Server
	if err := application.Run(config.GatewayPort); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
