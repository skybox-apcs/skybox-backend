package main

import (
	"sync"

	"skybox-backend/configs"
	"skybox-backend/docs" // Swagger documentation package
	apiApp "skybox-backend/internal/api/app"
	blockApp "skybox-backend/internal/blockserver/app"

	"go.uber.org/zap" // Third-party library for logging
)

func main() {
	// Load configs
	configs.LoadConfig()

	// Initialize Swagger
	docs.SwaggerInfo.Host = configs.Config.ServerHost + ":" + configs.Config.ServerPort

	var wg sync.WaitGroup            // WaitGroup to wait for all goroutines to finish
	logger, _ := zap.NewProduction() // Create a new logger instance
	defer logger.Sync()              // Flush any buffered log entries before exiting

	wg.Add(2) // Add two goroutines to the WaitGroup

	// Start the API server in a separate goroutine
	go func() {
		defer wg.Done() // Decrement the counter when the goroutine completes

		// Start the server
		apiApp.StartServer()
	}()

	// Start the Block server in a separate goroutine
	go func() {
		defer wg.Done() // Decrement the counter when the goroutine completes

		// Start the server
		blockApp.StartServer()
	}()

	wg.Wait() // Wait for all goroutines to finish
}
