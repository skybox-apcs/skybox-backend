// @title           Skybox API
// @version         1.0
// @description     Skybox is a cloud-based file storage provider similar to Google Drive and Dropbox. It allows users to securely store, manage, and retrieve their files.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Skybox Support
// @contact.url    http://skybox.io/support
// @contact.email  support@skybox.io

// @license.name  GNU General Public License v3.0
// @license.url   https://www.gnu.org/licenses/gpl-3.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in header
// @name Authorization
// @description  Use "Bearer {your_token}" to authenticate requests.

// @externalDocs.description  OpenAPI Documentation
// @externalDocs.url          https://swagger.io/resources/open-api/

package main

import (
	"skybox-backend/configs"
	"skybox-backend/docs"
	"skybox-backend/internal/api/app"
)

func main() {
	// Load Configurations
	configs.LoadConfig()

	// Initialize Swagger
	docs.SwaggerInfo.Host = configs.Config.ServerHost + ":" + configs.Config.ServerPort

	// Create a new application
	application := app.NewApplication()

	// Setup the DB
	db := application.Mongo.Database(configs.Config.MongoDBName)
	defer application.CloseDBConnection()

	// Start the server
	ginServer := app.NewServer()
	ginServer.CorsMiddleware()
	ginServer.SecurityMiddleware()
	ginServer.RouteMiddleware(db)
	ginServer.GlobalErrorHandler()

	ginServer.StartServer()
}
