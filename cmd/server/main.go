// @title           Skybox API and Block Server
// @version         1.0
// @description     Skybox is a cloud-based file storage provider similar to Google Drive and Dropbox. It allows users to securely store, manage, and retrieve their files.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Skybox Support
// @contact.url    None
// @contact.email  None

// @license.name  GNU General Public License v3.0
// @license.url   https://www.gnu.org/licenses/gpl-3.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description  Use "Bearer {your_token}" to authenticate requests.

// @externalDocs.description  OpenAPI Documentation
// @externalDocs.url          https://swagger.io/resources/open-api/

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
