package repositories

import (
	"context"
	"skybox-backend/internal/api/models"
	"slices"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type uploadSessionRepository struct {
	database   *mongo.Database
	collection string
}

func NewUploadSessionRepository(db *mongo.Database, collection string) *uploadSessionRepository {
	return &uploadSessionRepository{
		database:   db,
		collection: collection,
	}
}

// CreateSessionRecord creates a upload session record in the database
func (ur *uploadSessionRepository) CreateSessionRecord(ctx context.Context, session *models.UploadSession) (*models.UploadSession, error) {
	collection := ur.database.Collection(ur.collection)

	// Create the upload session in the database
	result, err := collection.InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}

	session.ID = result.InsertedID.(primitive.ObjectID)
	return session, nil
}

// GetSessionRecord retrieves a upload session record by session token
func (ur *uploadSessionRepository) GetSessionRecord(ctx context.Context, sessionToken string) (*models.UploadSession, error) {
	collection := ur.database.Collection(ur.collection)

	var session models.UploadSession
	err := collection.FindOne(ctx, bson.M{"session_token": sessionToken}).Decode(&session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetSessionRecordByFileID retrieves a upload session record by file ID
func (ur *uploadSessionRepository) GetSessionRecordByFileID(ctx context.Context, fileID string) (*models.UploadSession, error) {
	collection := ur.database.Collection(ur.collection)

	var session models.UploadSession
	err := collection.FindOne(ctx, bson.M{"file_id": fileID}).Decode(&session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// AddChunkSessionRecord adds a chunk to an existing upload session record
func (ur *uploadSessionRepository) AddChunkSessionRecord(ctx context.Context, sessionToken string, chunkNumber int) error {
	collection := ur.database.Collection(ur.collection)

	// Check if the chunk number already exists in the session record
	session, err := ur.GetSessionRecord(ctx, sessionToken)
	if err != nil {
		return err
	}
	if slices.Contains(session.ChunkList, chunkNumber) {
		return nil // Chunk already exists, no need to add it again
	}

	// Update the session record to add the chunk number
	_, err = collection.UpdateOne(ctx, bson.M{"session_token": sessionToken}, bson.M{
		"$addToSet": bson.M{"chunk_list": chunkNumber},
	})
	if err != nil {
		return err
	}

	// Update the session record if the chunk number is the last chunk
	if len(session.ChunkList)+1 == session.TotalChunks {
		// Update the session record to mark it as completed
		_, err = collection.UpdateOne(ctx, bson.M{"session_token": sessionToken}, bson.M{
			"$set": bson.M{"status": "completed"},
		})
		if err != nil {
			return err
		}

		// Update the file record to mark it as completed
		fileCollection := ur.database.Collection(models.CollectionFiles)
		_, err = fileCollection.UpdateOne(ctx, bson.M{"_id": session.FileID}, bson.M{
			"$set": bson.M{"status": "uploaded"},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
