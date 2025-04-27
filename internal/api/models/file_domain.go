package models

type RenameFileRequest struct {
	NewName string `json:"new_name" binding:"required"`
}

type RenameFileResponse struct {
}

type MoveFileRequest struct {
	NewParentID string `json:"new_parent_id" binding:"required"`
}

type MoveFileResponse struct {
}
