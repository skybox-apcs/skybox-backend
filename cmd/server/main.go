// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

package main

import (
	"skybox-backend/configs"
	_ "skybox-backend/docs"
	"skybox-backend/internal/app"
)

func main() {
	// Load Configurations
	configs.LoadConfig()

	// Create a new application
	application := app.NewApplication()

	// Setup the DB
	db := application.Mongo.Database(configs.Config.MongoDBName)
	defer application.CloseDBConnection()

	// Start the server
	ginServer := app.NewServer()
	ginServer.SecurityMiddleware()
	ginServer.RouteMiddleware(db)
	ginServer.GlobalErrorHandler()

	ginServer.StartServer()
}
