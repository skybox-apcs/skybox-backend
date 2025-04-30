package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"skybox-backend/configs"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/services"
	"skybox-backend/internal/shared"

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

func parseRangeHeader(rangeHeader string, fileSize int64) (int64, int64) {
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return 0, fileSize - 1 // default to the whole file
	}

	parts := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
	start, _ := strconv.ParseInt(parts[0], 10, 64)
	var end int64
	if len(parts) > 1 && parts[1] != "" {
		end, _ = strconv.ParseInt(parts[1], 10, 64)
	} else {
		end = fileSize - 1 // default to the end of the file
	}

	if end >= fileSize {
		end = fileSize - 1 // ensure end is within bounds
	}

	return start, end
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

// PartialDownloadFileHandler godoc
func (fc *FileController) PartialDownloadFileHandler(c *gin.Context) {
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

	// Check if the Range header is present
	rangeHeader := c.GetHeader("Range")
	start, end := parseRangeHeader(rangeHeader, file.Size)
	if start > end || end >= file.Size {
		shared.RespondJson(c, http.StatusRequestedRangeNotSatisfiable, "error", "Invalid range", nil)
		return
	}

	// Prepare buffer for the request range
	buf := new(bytes.Buffer)
	chunkSize := configs.Config.DefaultChunkSize
	startChunk := start / chunkSize               // round down to the previous chunk
	endChunk := (end + chunkSize - 1) / chunkSize // round up to the next chunk

	// Iterate over the chunks and download them
	for i := startChunk; i <= endChunk; i++ {
		chunkData, err := fc.FileService.GetFileData(c, fileID, int(i))
		if err != nil {
			c.Error(err)
			return
		}

		chunkStart := int64(i) * int64(chunkSize)
		forwardStart := max(0, start-chunkStart)
		forwardEnd := min(int64(len(chunkData)), end-chunkStart+1)

		buf.Write(chunkData[forwardStart:forwardEnd])
	}

	// Set the headers for the response
	c.Header("Content-Type", file.MimeType)
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, file.Size))
	c.Status(http.StatusPartialContent)
	c.Writer.Write(buf.Bytes())
}

// FullDownloadFileHandler godoc
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

	// Set the headers for the response
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.FileName))
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Length", strconv.FormatInt(file.Size, 10))
	c.Header("Accept-Ranges", "bytes")
	c.Status(http.StatusOK)

	// Stream the file data to the response
	for i := 0; i < file.TotalChunks; i++ {
		chunkData, err := fc.FileService.GetFileData(c, fileID, i)
		if err != nil {
			c.Error(err)
			return
		}

		// Stream the chunk data to the response
		if _, err := c.Writer.Write(chunkData); err != nil {
			c.Error(err)
			return
		}
		c.Writer.Flush()
	}
}
