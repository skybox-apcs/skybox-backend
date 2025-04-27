package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionFolders = "folders"
)

type FolderStat struct {
	TotalFiles   int   `bson:"total_files" json:"total_files"`
	TotalFolders int   `bson:"total_folders" json:"total_folders"`
	TotalSize    int64 `bson:"total_size" json:"total_size"`
}

// Folder struct encapsulates the folder model
type Folder struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OwnerID        primitive.ObjectID `bson:"owner_id" json:"owner_id"`                                     // The owner of the folder
	ParentFolderID primitive.ObjectID `bson:"parent_folder_id,omitempty" json:"parent_folder_id,omitempty"` // The parent folder ID, if any
	Name           string             `bson:"name" json:"name"`
	IsDeleted      bool               `bson:"is_deleted" json:"is_deleted"`
	Stats          FolderStat         `bson:"stats" json:"stats"`
	IsRoot         bool               `bson:"is_root" json:"is_root"` // Indicates if this is a root folder

	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"` // Nullable field for soft delete
}

type FolderRepository interface {
	CreateFolder(ctx context.Context, folder *Folder) (*Folder, error)
	GetFolderByID(ctx context.Context, id string) (*Folder, error)
	GetFolderParentIDByFolderID(ctx context.Context, folderID string) (string, error)
	GetFolderListInFolder(ctx context.Context, folderID string) ([]*Folder, error)
	GetFileListInFolder(ctx context.Context, folderID string) ([]*File, error)
	DeleteFolder(ctx context.Context, id string) error
	RenameFolder(ctx context.Context, id string, newName string) error
	MoveFolder(ctx context.Context, id string, newParentID string) error
}
