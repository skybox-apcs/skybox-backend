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
	Status    string     `bson:"status" json:"status"`                             // Status of the file (e.g., "uploaded", "processing", "failed")

	// ChunkList   []Chunk `bson:"chunk_list" json:"chunk_list"`
	TotalChunks int `bson:"total_chunks" json:"total_chunks"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`

	OwnerEmail    string `bson:"owner_email,omitempty" json:"owner_email,omitempty"`
	OwnerUsername string `bson:"owner_username,omitempty" json:"owner_username,omitempty"`
}

type FileRepository interface {
	UploadFileMetadata(ctx context.Context, file *File) (*File, error) // Upload file metadata
	GetFileByID(ctx context.Context, id string) (*File, error)         // Get FilenewParentID metadata
	DeleteFile(ctx context.Context, id string) error
	RenameFile(ctx context.Context, id string, newName string) error
	MoveFile(ctx context.Context, id string, newParentFolderID string) error
	SearchFiles(ctx context.Context, ownerId primitive.ObjectID, query string) ([]*File, error) // Search files by name or folder name
}
