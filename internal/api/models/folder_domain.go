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
	Contents []any `json:"contents"`
}
