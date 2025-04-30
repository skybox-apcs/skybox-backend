package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	CollectionChunks = "chunks"
)

type Chunk struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FileID     primitive.ObjectID `bson:"file_id" json:"file_id"`
	ChunkIndex int                `bson:"chunk_index" json:"chunk_index"`
	ChunkSize  int64              `bson:"chunk_size" json:"chunk_size"` // Size of the chunk in bytes
	ChunkHash  string             `bson:"chunk_hash" json:"chunk_hash"` // Hash of the chunk data, could use MD5, SHA-1, etc.
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

type ChunkRepository interface {
	UploadChunkMetadata(fileId string, chunk *Chunk) (*Chunk, error)      // Upload chunk metadata
	UpdateChunkStatus(fileId string, chunkIndex int, status string) error // Update chunk status
	GetChunksByFileID(fileId string) ([]Chunk, error)                     // Get all chunks for a file
	GetChunkByID(id string) (*Chunk, error)                               // Get chunk by ID
}
