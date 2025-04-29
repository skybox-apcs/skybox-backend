package app

import (
	"context"
	"fmt"
	"time"

	"skybox-backend/configs"

	"go.mongodb.org/mongo-driver/bson"
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

func CreateIndexes(db *mongo.Database) error {
	// Create a context with a timeout for the index creation
	// This is important to avoid blocking the application indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define a string - index mapping for the indexes to be created
	indexes := map[string][]mongo.IndexModel{}

	// Define the indexes for the "users" collection
	indexes["users"] = []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}}, // Unique index on email
			Options: options.Index().SetUnique(true),
		},
	}

	// Define the indexes for the "files" collection
	indexes["files"] = []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "owner_id", Value: 1},         // Index on owner_id
				{Key: "parent_folder_id", Value: 1}, // Index on parent_folder_id
				{Key: "name", Value: -1},            // Sort by name descending
			}, // Sort by created_at descending
		},
	}

	// Define the indexes for the "folders" collection
	indexes["folders"] = []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "owner_id", Value: 1},         // Index on owner_id
				{Key: "parent_folder_id", Value: 1}, // Index on parent_folder_id
				{Key: "name", Value: -1},            // Sort by name descending
			}, // Sort by created_at descending
		},
	}

	// Define the indexes for the "user_tokens" collection
	indexes["user_tokens"] = []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1}, // Index on user_id
				{Key: "token", Value: 1},   // Index on token
			},
		},
	}

	// Create the indexes for each collection using goroutines
	for collectionName, indexModels := range indexes {
		collection := db.Collection(collectionName)
		_, err := collection.Indexes().CreateMany(ctx, indexModels)
		if err != nil {
			return fmt.Errorf("failed to create indexes for collection %s: %v", collectionName, err)
		}
		fmt.Printf("Indexes created for collection %s\n", collectionName)
	}

	return nil
}
