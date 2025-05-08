package controllers

import (
	"fmt"
	"net/http"

	"skybox-backend/configs"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/services"
	"skybox-backend/internal/shared"
	"skybox-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// FileController handles file-related requests
type FileController struct {
	FileService *services.FileService
}

// NewFileController creates a new instance of FileController
func NewFileController(fileService *services.FileService) *FileController {
	return &FileController{
		FileService: fileService,
	}
}

// getFileMetadataHandler godoc
//
// @Summary Get file metadata
// @Description Get file metadata by file ID. The file ID is a unique identifier for the file in the database. The metadata includes information such as the file name, size, type, and creation date.
// @Security		Bearer
// @Tags Files
// @Accept json
// @Produce json
// @Param fileId path string true "File ID" example(1234567890abcdef12345678)
// @Success 200 {object} models.File "File metadata retrieved successfully"
// @Failure 400 {string} string Invalid request body"
// @Failure 404 {string} string "File not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/files/{fileId} [get]
func (fc *FileController) GetFileMetadataHandler(c *gin.Context) {
	// Get the file ID from the URL parameters
	fileID := c.Param("fileId")
	if fileID == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "File ID is required", nil)
		return
	}

	file, err := fc.FileService.GetFileByID(c, fileID)
	if err != nil {
		c.Error(err)
		return
	}

	// Send the response
	shared.RespondJson(c, http.StatusOK, "success", "File metadata retrieved successfully", file)
}

// deleteFileHandler godoc
//
// @Summary Delete a file
// @Description Delete a file by its ID. The file is marked as deleted in the database, but not removed from the storage service. This allows for potential recovery in the future.
// @Security		Bearer
// @Tags Files
// @Accept json
// @Produce json
// @Param fileId path string true "File ID" example(1234567890abcdef12345678)
// @Success 200 {string} string "File deleted successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "File not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/files/{fileId} [delete]
func (fc *FileController) DeleteFileHandler(c *gin.Context) {
	// Get the file ID from the URL parameters
	fileID := c.Param("fileId")
	if fileID == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "File ID is required", nil)
		return
	}

	err := fc.FileService.DeleteFile(c, fileID)
	if err != nil {
		c.Error(err)
		return
	}

	// Send the response
	shared.RespondJson(c, http.StatusOK, "success", "File deleted successfully", nil)
}

// RenameFileHandler godoc
//
// @Summary Rename a file
// @Description Rename a file by providing the new name. Upon renaming, the file's metadata is updated. And when downloading, the file is fetched from the storage service using the new name.
// @Security		Bearer
// @Tags Files
// @Accept json
// @Produce json
// @Param fileId path string true "File ID" example(1234567890abcdef12345678)
// @Param request body models.RenameFileRequest true "Rename file request"
// @Success 200 {string} string "File renamed successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "File not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/files/{fileId}/rename [put]
// @Router /api/v1/files/{fileId}/rename [patch]
func (fc *FileController) RenameFileHandler(c *gin.Context) {
	// Get the file ID from the URL parameters
	fileID := c.Param("fileId")
	if fileID == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "File ID is required", nil)
		return
	}

	// Get the new name from the request body
	var requestBody models.RenameFileRequest
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	// Rename the file using the service
	err = fc.FileService.RenameFile(c, fileID, requestBody.NewName)
	if err != nil {
		c.Error(err)
		return
	}

	// Send the response
	shared.RespondJson(c, http.StatusOK, "success", "File renamed successfully", nil)
}

// MoveFileHandler godoc
//
// @Summary Move a file to a new folder
// @Description Move a file to a new folder by providing the new parent folder ID.
// @Security		Bearer
// @Tags Files
// @Accept json
// @Produce json
// @Param fileId path string true "File ID" example(1234567890abcdef12345678)
// @Param request body models.MoveFileRequest true "Move file request"
// @Success 200 {string} string "File moved successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "File not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/files/{fileId}/move [put]
func (fc *FileController) MoveFileHandler(c *gin.Context) {
	// Get the file ID from the URL parameters
	fileID := c.Param("fileId")
	if fileID == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "File ID is required", nil)
		return
	}

	// Get the new parent folder ID from the request body
	var requestBody models.MoveFileRequest
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		shared.RespondJson(c, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	// Move the file using the service
	err = fc.FileService.MoveFile(c, fileID, requestBody.NewParentID)
	if err != nil {
		c.Error(err)
		return
	}

	// Send the response
	shared.RespondJson(c, http.StatusOK, "success", "File moved successfully", nil)
}

// FullDownloadFileHandler godoc
//
// @Summary Download a file
// @Description Download a file by its ID. The file is streamed in chunks to the client. The client can request a specific range of bytes using the Range header. If no Range header is provided, the entire file is downloaded.
// @Security		Bearer
// @Tags Files
// @Accept json
// @Produce json
// @Param fileId path string true "File ID" example(1234567890abcdef12345678)
// @Success 200 {string} string "File downloaded successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "File not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/files/{fileId}/download [get]
func (fc *FileController) FullDownloadFileHandler(c *gin.Context) {
	// Get the file ID from the URL parameters
	fileID := c.Param("fileId")
	if fileID == "" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "File ID is required", nil)
		return
	}

	// Get the file metadata
	file, err := fc.FileService.GetFileByID(c, fileID)
	if err != nil {
		c.Error(err)
		return
	}

	// Check if the file uploaded
	if file.Status != "uploaded" {
		shared.RespondJson(c, http.StatusBadRequest, "error", "File is not uploaded yet", nil)
		return
	}

	// Generate token
	token, err := utils.GenerateToken(
		map[string]string{
			"fileId":      fileID,
			"ownerId":     file.OwnerID.Hex(),
			"totalChunks": fmt.Sprintf("%d", file.TotalChunks),
			"fileName":    file.FileName,
			"fileSize":    fmt.Sprintf("%d", file.Size),
		},
		configs.Config.JWTSecret,
		1,
	)
	if err != nil {
		c.Error(err)
		return
	}

	// Generate the download URL for the block server
	downloadURL := fmt.Sprintf("http://%s:%s/download/%s?token=%s",
		configs.Config.BlockServerHost,
		configs.Config.BlockServerPort,
		fileID,
		token,
	)
	c.Redirect(http.StatusFound, downloadURL)
}
