package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	CollectionUploadSessions = "upload_sessions"
)

type UploadSession struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`             // Reference to the user
	FileID       primitive.ObjectID `bson:"file_id" json:"file_id"`             // Reference to the file
	SessionToken string             `bson:"session_token" json:"session_token"` // Unique session token for the upload session
	TotalSize    int64              `bson:"total_size" json:"total_size"`       // Total size of the file to be uploaded
	ActualSize   int64              `bson:"actual_size" json:"actual_size"`     // Actual size of the uploaded file
	ChunkList    []int              `bson:"chunk_list" json:"chunk_list"`       // List of chunks that have been uploaded
	Status       string             `bson:"status" json:"status"`               // Status of the upload session (e.g., "pending", "completed", "failed")
}

type UploadSessionRepository interface {
	CreateSessionRecord(ctx context.Context, session *UploadSession) (*UploadSession, error)
	GetSessionRecord(ctx context.Context, sessionToken string) (*UploadSession, error)
	GetSessionRecordByFileID(ctx context.Context, fileID string) (*UploadSession, error)
	AddChunkSessionRecord(ctx context.Context, sessionToken string, chunkNumber int, chunkSize int) error
	AddChunkSessionRecordByFileID(ctx context.Context, fileID string, chunkNumber int, chunkSize int) error
}
