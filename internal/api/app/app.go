package app

import (
	"fmt"

	"skybox-backend/configs"

	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	Mongo *mongo.Client
}

func NewApplication() Application {
	app := &Application{}

	// Connect to MongoDB
	app.Mongo = NewMongoDatabase()
	fmt.Println("Connected to MongoDB")

	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDatabase(app.Mongo)
	fmt.Println("Closed MongoDB connection")
}

func StartServer() {
	// Create a new application
	application := NewApplication()
	defer application.CloseDBConnection()

	// Setup the DB
	db := application.Mongo.Database(configs.Config.MongoDBName)

	// Setup the DB Indexes
	CreateIndexes(db)
	if err := CreateIndexes(db); err != nil {
		panic(err)
	}

	// Start the server
	ginServer := NewServer()

	// Setup middleware in order of execution
	ginServer.CorsMiddleware()
	ginServer.GlobalErrorHandler()
	ginServer.SecurityMiddleware()
	ginServer.RateLimitMiddleware()
	ginServer.RouteMiddleware(db)

	// Start the server
	ginServer.StartServer()
}
