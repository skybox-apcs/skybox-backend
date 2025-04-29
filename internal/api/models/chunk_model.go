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
	Status     string             `bson:"status" json:"status"` // "pending", "uploaded", "failed"
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

type ChunkRepository interface {
	UploadChunkMetadata(fileId string, chunk *Chunk) (*Chunk, error)      // Upload chunk metadata
	UpdateChunkStatus(fileId string, chunkIndex int, status string) error // Update chunk status
	GetChunksByFileID(fileId string) ([]Chunk, error)                     // Get all chunks for a file
	GetChunkByID(id string) (*Chunk, error)                               // Get chunk by ID
}
