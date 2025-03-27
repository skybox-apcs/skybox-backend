package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient holds the MongoDB client and database instances
type MongoClient struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// NewMongoClient creates a new MongoDB client
func NewMongoClient(URI string, dbName string) (*MongoClient, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(URI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &MongoClient{
		Client: client,
		DB:     client.Database(dbName),
	}, nil
}
