package models

import "time"

type CreateFolderRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateFolderResponse struct {
	ID             string `json:"id"`
	OwnerID        string `json:"owner_id"`
	ParentFolderID string `json:"parent_id"`
	Name           string `json:"name"`
}

type GetFolderContentsRequest struct{}

type GetFolderContentsResponse struct {
	FolderList []*FolderResponse `json:"folder_list"`
	FileList   []*FileResponse   `json:"file_list"`
}

type RenameFolderRequest struct {
	NewName string `json:"new_name" binding:"required"`
}

type RenameFolderResponse struct {
}

type MoveFolderRequest struct {
	NewParentID string `json:"new_parent_id" binding:"required"`
}

type MoveFolderResponse struct {
}

type FileResponse struct {
	ID             string    `json:"id" bson:"_id, omitempty"`
	ParentFolderID string    `json:"parent_folder_id" bson:"parent_folder_id"`
	OwnerID        string    `json:"owner_id" bson:"owner_id"`
	OwnerUsername  string    `json:"owner_user_name" bson:"owner_user_name"`
	OwnerEmail     string    `json:"owner_email" bson:"owner_email"`
	Name           string    `json:"name" bson:"name"`
	MimeType       string    `json:"mime_type" bson:"mime_type"`
	Size           int64     `json:"size" bson:"size"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}

type FolderResponse struct {
	ID             string     `json:"id" bson:"_id, omitempty"`
	ParentFolderID string     `json:"parent_folder_id" bson:"parent_folder_id"`
	OwnerID        string     `json:"owner_id" bson:"owner_id"`
	OwnerUsername  string     `json:"owner_user_name" bson:"owner_user_name"`
	OwnerEmail     string     `json:"owner_email" bson:"owner_email"`
	Name           string     `json:"name" bson:"name"`
	Stats          FolderStat `json:"stats" bson:"stats"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" bson:"updated_at"`
}

type UploadFileMetadataRequest struct {
	FileName string `json:"file_name" binding:"required"`
	FileSize int64  `json:"file_size" binding:"required"`
	MimeType string `json:"mime_type"` // optional
}

type UploadFileMetadataResponse struct {
	File      FileResponse `json:"file"`
	UploadURL string       `json:"upload_url" binding:"required"`
}

type UpdateFolderPublicRequest struct {
	IsPublic bool `json:"is_public"`
}

type ShareFolderRequest struct {
	UserID     string `json:"user_id" binding:"required"`
	Permission bool   `json:"permission"`
}

type RevokeFolderShareRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type RemoveFolderShareRequest struct {
	UserID string `json:"user_id" binding:"required"`
}
