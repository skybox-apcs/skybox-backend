package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionFiles = "files"
)

// File struct encapsulates the file model
type File struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OwnerID        primitive.ObjectID `bson:"owner_id" json:"owner_id"`                                     // The owner of the file
	ParentFolderID primitive.ObjectID `bson:"parent_folder_id,omitempty" json:"parent_folder_id,omitempty"` // The parent folder ID, if any

	FileName  string     `bson:"file_name" json:"file_name"`
	MimeType  string     `bson:"mime_type" json:"mime_type"`
	Extension string     `bson:"extension" json:"extension"`
	Size      int64      `bson:"size" json:"size"`
	IsDeleted bool       `bson:"is_deleted" json:"is_deleted"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"` // Nullable field for soft delete

	// ChunkList   []Chunk `bson:"chunk_list" json:"chunk_list"`
	TotalChunks int `bson:"total_chunks" json:"total_chunks"`
}

type FileRepository interface {
	CreateFile(ctx context.Context, file *File) error
	GetFileByID(ctx context.Context, id string) (*File, error)
	DeleteFile(ctx context.Context, id string) error
}
