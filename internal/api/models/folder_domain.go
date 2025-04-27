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
	FolderList []*Folder `json:"folder_list"`
	FileList   []*File   `json:"file_list"`
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
	ID             string    `json:"id"`
	OwnerID        string    `json:"owner_id"`
	ParentFolderID string    `json:"parent_id"`
	Name           string    `json:"name"`
	MimeType       string    `json:"mime_type"`
	Size           int64     `json:"size"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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
