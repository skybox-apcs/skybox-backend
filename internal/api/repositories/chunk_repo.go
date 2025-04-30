package repositories

import (
	"context"
	"skybox-backend/internal/api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChunkRepository struct {
	database   *mongo.Database
	collection string
}

// NewChunkRepository creates a new instance of the ChunkRepository
func NewChunkRepository(db *mongo.Database, collection string) *ChunkRepository {
	return &ChunkRepository{
		database:   db,
		collection: collection,
	}
}

// UploadChunkMetadata uploads chunk metadata to the database
func (cr *ChunkRepository) UploadChunkMetadata(ctx context.Context, fileId string, chunk *models.Chunk) (*models.Chunk, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Insert the chunk metadata into the database
	result, err := cr.database.Collection(cr.collection).InsertOne(ctx, chunk)
	if err != nil {
		return nil, err
	}

	// Set the ID of the chunk to the inserted ID
	chunk.ID = result.InsertedID.(primitive.ObjectID)

	return chunk, nil
}

// UpdateChunkStatus updates the status of a chunk in the database
func (cr *ChunkRepository) UpdateChunkStatus(ctx context.Context, fileId string, chunkIndex int, status string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Update the chunk status in the database
	_, err := cr.database.Collection(cr.collection).UpdateOne(ctx, bson.M{
		"file_id":     fileId,
		"chunk_index": chunkIndex,
	}, bson.M{
		"$set": bson.M{
			"status": status,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// GetChunksByFileID retrieves all chunks for a specific file from the database
func (cr *ChunkRepository) GetChunksByFileID(ctx context.Context, fileId string) ([]models.Chunk, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Retrieve all chunks for the specified file ID
	cursor, err := cr.database.Collection(cr.collection).Find(ctx, bson.M{"file_id": fileId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chunks []models.Chunk
	if err := cursor.All(ctx, &chunks); err != nil {
		return nil, err
	}

	return chunks, nil
}

// GetChunkByID retrieves a chunk by its ID from the database
func (cr *ChunkRepository) GetChunkByID(ctx context.Context, id string) (*models.Chunk, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Convert the ID string to an ObjectID
	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Retrieve the chunk by ID
	var chunk models.Chunk
	err = cr.database.Collection(cr.collection).FindOne(ctx, bson.M{"_id": idHex}).Decode(&chunk)
	if err != nil {
		return nil, err
	}

	return &chunk, nil
}
