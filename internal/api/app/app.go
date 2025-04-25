package app

import (
	"skybox-backend/configs"

	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	Mongo *mongo.Client
}

func NewApplication() Application {
	app := &Application{}

	app.Mongo = NewMongoDatabase() // Connect to the MongoDB

	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDatabase(app.Mongo)
}

func StartServer() {
	// Create a new application
	application := NewApplication()

	// Setup the DB
	db := application.Mongo.Database(configs.Config.MongoDBName)
	defer application.CloseDBConnection()

	// Start the server
	ginServer := NewServer()
	ginServer.CorsMiddleware()
	ginServer.SecurityMiddleware()
	ginServer.RouteMiddleware(db)
	ginServer.GlobalErrorHandler()

	ginServer.StartServer()
}
