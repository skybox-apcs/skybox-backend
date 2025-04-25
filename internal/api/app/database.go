package app

import (
	"context"
	"fmt"

	"skybox-backend/configs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDatabase() *mongo.Client {
	// Set the options
	opts := options.Client().ApplyURI(configs.Config.MongoURI)

	// Connect to the MongoDB
	fmt.Println("Connecting to MongoDB...")
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB")
	return client
}

// CloseMongoDatabase closes the MongoDB connection
func CloseMongoDatabase(client *mongo.Client) {
	err := client.Disconnect(context.Background())
	if err != nil {
		panic(err)
	}
}
