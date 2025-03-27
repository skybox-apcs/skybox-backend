package main

import (
	"skybox-backend/configs"

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
