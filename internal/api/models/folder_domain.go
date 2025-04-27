package models

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
