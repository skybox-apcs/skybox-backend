package app

import (
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
