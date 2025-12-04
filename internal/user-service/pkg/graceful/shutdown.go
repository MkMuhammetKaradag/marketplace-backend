package graceful

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func WaitForShutdown(app *fiber.App, timeout time.Duration, ctx context.Context) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	fmt.Println("Shutdown signal received")

	if err := app.ShutdownWithTimeout(timeout); err != nil {
		fmt.Println("Error during server shutdown:", err)
	}

	fmt.Println("Server gracefully stopped")
}
