package repositories

import (
	"context"
	"skybox-backend/internal/api/models"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UploadSessionRepository struct {
	database   *mongo.Database
	collection string
}

func NewUploadSessionRepository(db *mongo.Database, collection string) *UploadSessionRepository {
	return &UploadSessionRepository{
		database:   db,
		collection: collection,
	}
}

// CreateSessionRecord creates a upload session record in the database
func (ur *UploadSessionRepository) CreateSessionRecord(ctx context.Context, session *models.UploadSession) (*models.UploadSession, error) {
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
func (ur *UploadSessionRepository) GetSessionRecord(ctx context.Context, sessionToken string) (*models.UploadSession, error) {
	collection := ur.database.Collection(ur.collection)

	var session models.UploadSession
	err := collection.FindOne(ctx, bson.M{"session_token": sessionToken}).Decode(&session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetSessionRecordByFileID retrieves a upload session record by file ID
func (ur *UploadSessionRepository) GetSessionRecordByFileID(ctx context.Context, fileID string) (*models.UploadSession, error) {
	collection := ur.database.Collection(ur.collection)

	var session models.UploadSession
	fileIDObj, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(ctx, bson.M{"file_id": fileIDObj}).Decode(&session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetSessionRecordByUserID retrieves upload session records by user ID
func (ur *UploadSessionRepository) GetSessionRecordByUserID(ctx context.Context, userID string) (*[]models.UploadSession, error) {
	collection := ur.database.Collection(ur.collection)

	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var sessions []models.UploadSession
	cursor, err := collection.Find(ctx, bson.M{"user_id": userIDObj})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &sessions, nil // No sessions found, return empty slice
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}

	return &sessions, nil
}

func (ur *UploadSessionRepository) AddChunkSessionRecord(ctx context.Context, sessionToken string, chunkNumber int, chunkSize int, chunkHash string) error {
	session, err := ur.database.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// Run all DB operations in a transaction
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		collection := ur.database.Collection(ur.collection)
		chunkCollection := ur.database.Collection(models.CollectionChunks)
		fileCollection := ur.database.Collection(models.CollectionFiles)

		// Re-fetch inside transaction
		sessionRecord, err := ur.GetSessionRecord(sessCtx, sessionToken)
		if err != nil {
			return nil, err
		}
		if slices.Contains(sessionRecord.ChunkList, chunkNumber) {
			return nil, nil // Already added
		}

		// Insert chunk
		if _, err := chunkCollection.InsertOne(sessCtx, bson.M{
			"file_id":     sessionRecord.FileID,
			"chunk_index": chunkNumber,
			"chunk_size":  int64(chunkSize),
			"chunk_hash":  chunkHash,
			"created_at":  time.Now(),
		}); err != nil {
			return nil, err
		}

		// Atomically update session
		update := bson.M{
			"$addToSet": bson.M{"chunk_list": chunkNumber},
			"$inc":      bson.M{"actual_size": int64(chunkSize)},
		}
		if _, err := collection.UpdateOne(sessCtx, bson.M{"session_token": sessionToken}, update); err != nil {
			return nil, err
		}

		// Re-fetch updated session to decide if it's complete
		sessionRecord, err = ur.GetSessionRecord(sessCtx, sessionToken)
		if err != nil {
			return nil, err
		}
		if sessionRecord.TotalSize <= sessionRecord.ActualSize {
			// Complete session and file
			if _, err := collection.UpdateOne(sessCtx, bson.M{"session_token": sessionToken}, bson.M{
				"$set": bson.M{"status": "completed"},
			}); err != nil {
				return nil, err
			}
			if _, err := fileCollection.UpdateOne(sessCtx, bson.M{"_id": sessionRecord.FileID}, bson.M{
				"$set": bson.M{
					"status":       "uploaded",
					"total_chunks": len(sessionRecord.ChunkList) + 1, // +1 for current chunk
				},
			}); err != nil {
				return nil, err
			}
		}

		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}

// AddChunkSessionRecordByFileID adds a chunk to an existing upload session record using file ID
func (ur *UploadSessionRepository) AddChunkSessionRecordByFileID(ctx context.Context, fileID string, chunkNumber int, chunkSize int, chunkHash string) error {
	session, err := ur.database.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		collection := ur.database.Collection(ur.collection)
		chunkCollection := ur.database.Collection(models.CollectionChunks)
		fileCollection := ur.database.Collection(models.CollectionFiles)

		// Fetch session inside transaction
		sessionRecord, err := ur.GetSessionRecordByFileID(sessCtx, fileID)
		if err != nil {
			return nil, err
		}
		if slices.Contains(sessionRecord.ChunkList, chunkNumber) {
			return nil, nil // Already uploaded
		}

		// Insert chunk
		if _, err := chunkCollection.InsertOne(sessCtx, bson.M{
			"file_id":     sessionRecord.FileID,
			"chunk_index": chunkNumber,
			"chunk_size":  int64(chunkSize),
			"chunk_hash":  chunkHash,
			"created_at":  time.Now(),
		}); err != nil {
			return nil, err
		}

		// Update session
		if _, err := collection.UpdateOne(sessCtx, bson.M{"file_id": sessionRecord.FileID}, bson.M{
			"$addToSet": bson.M{"chunk_list": chunkNumber},
			"$inc":      bson.M{"actual_size": int64(chunkSize)},
		}); err != nil {
			return nil, err
		}

		// Re-fetch updated session
		sessionRecord, err = ur.GetSessionRecordByFileID(sessCtx, fileID)
		if err != nil {
			return nil, err
		}
		if sessionRecord.TotalSize <= sessionRecord.ActualSize {
			if _, err := collection.UpdateOne(sessCtx, bson.M{"file_id": sessionRecord.FileID}, bson.M{
				"$set": bson.M{"status": "completed"},
			}); err != nil {
				return nil, err
			}
			if _, err := fileCollection.UpdateOne(sessCtx, bson.M{"_id": sessionRecord.FileID}, bson.M{
				"$set": bson.M{
					"status":       "uploaded",
					"total_chunks": len(sessionRecord.ChunkList) + 1, // include current
				},
			}); err != nil {
				return nil, err
			}
		}

		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}
