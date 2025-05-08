package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"skybox-backend/configs"
	"skybox-backend/internal/blockserver/services"
	"skybox-backend/internal/shared"
	"skybox-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// DownloadController handles file download requests
type DownloadController struct {
	downloadService *services.DownloadService
}

// NewDownloadController creates a new instance of DownloadController
func NewDownloadController(downloadService *services.DownloadService) *DownloadController {
	return &DownloadController{
		downloadService: downloadService,
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

// DownloadFileHandler godoc
//
// @Summary Download a file
// @Description Download a file by its ID. The file ID is a unique identifier for the file in the database. The file is downloaded in chunks to optimize performance and reduce memory usage.
// @Tags Files
// @Accept json
// @Produce json
// @Param fileId path string true "File ID" example(1234567890abcdef12345678)
// @Param token query string true "Token" example(eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaWxlTmFtZSI6InRlc3QuanBnIiwib3duZXJJZCI6IjEyMzQ1Njc4OWFiY2RlZiIsImZpbGVTaXplIjoxMjM0NTY3ODkwLCJ0b3RhbENodW5rcyI6MTIzNDU2Nzg5MCwibWltZVR5cGUiOiJpbWFnZS9qcGcifQ.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaWxlTmFtZSI6InRlc3QuanBnIiwib3duZXJJZCI6IjEyMzQ1Njc4OWFiY2RlZiIsImZpbGVTaXplIjoxMjM0NTY3ODkwLCJ0b3RhbENodW5rcyI6MTIzNDU2Nzg5MCwibWltZVR5cGUiOiJpbWFnZS9qcGcifQ)
// @Success 200 {string} string "File downloaded successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "File not found"
// @Failure 500 {string} string "Internal server error"
// @Router /download/{fileId} [get]
// DownloadFileHandler handles the file download request
func (dc *DownloadController) DownloadFileHandler(c *gin.Context) {
	// Since everything is vefied by the API Server, we don't need to verify the file metadata again
	fileID := c.Param("fileId")
	if fileID == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "File ID is required")
		return
	}

	// Get the token and validate
	token := c.Query("token")
	if token == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "Token is required")
		return
	}
	data, err := utils.GetKeysFromToken(token, configs.Config.JWTSecret)
	if err != nil {
		shared.ErrorJSON(c, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Retrieve data from the token
	fileName := data["fileName"]
	ownerId := data["ownerId"]
	fileSize, err := strconv.ParseInt(data["fileSize"], 10, 64)
	if err != nil {
		shared.ErrorJSON(c, http.StatusBadRequest, "Invalid file size")
		return
	}
	totalChunks, err := strconv.ParseInt(data["totalChunks"], 10, 64)
	if err != nil {
		shared.ErrorJSON(c, http.StatusBadRequest, "Invalid total chunks")
		return
	}
	mimeType := data["mimeType"]
	if mimeType == "" {
		mimeType = "application/octet-stream" // Default MIME type if not provided
	}

	// Set the headers for the response
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("Content-Type", mimeType)
	c.Header("Accept-Ranges", "bytes")

	// Parse the Range header
	rangeHeader := c.GetHeader("Range")
	if rangeHeader == "" {
		// If no Range header, send the full file
		c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
		c.Status(http.StatusOK)

		// Stream the file data to the response
		for i := 0; i < int(totalChunks); i++ {
			chunkData, err := dc.downloadService.DownloadFile(c, ownerId, fileID, i)
			if err != nil {
				shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to download file chunk")
				return
			}

			// Stream the chunk data to the response
			if _, err := c.Writer.Write(chunkData); err != nil {
				shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to write chunk data to response")
				return
			}
			c.Writer.Flush()
		}

		return
	}

	// If Range header is present, handle partial content
	start, end := parseRangeHeader(rangeHeader, fileSize)
	if start > end || end >= fileSize {
		shared.RespondJson(c, http.StatusRequestedRangeNotSatisfiable, "error", "Invalid range", nil)
		return
	}

	c.Header("Content-Length", strconv.FormatInt(end-start+1, 10))
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Status(http.StatusPartialContent)

	// Stream the file data to the response
	buf := new(bytes.Buffer)
	chunkSize := configs.Config.DefaultChunkSize
	startChunk := start / chunkSize             // round down to the previous chunk
	endChunk := (end+chunkSize-1)/chunkSize - 1 // round up to the next chunk and subtract 1 (to get the last chunk index)

	// DEBUG
	// fmt.Printf("ChunkSize: %d, start: %d, end: %d, startChunk: %d, endChunk: %d\n", chunkSize, start, end, startChunk, endChunk)

	// Iterate over the chunks and download them
	for i := startChunk; i <= endChunk; i++ {
		chunkData, err := dc.downloadService.DownloadFile(c, ownerId, fileID, int(i))
		if err != nil {
			c.Error(err)
			return
		}

		chunkStart := int64(i) * int64(chunkSize)
		forwardStart := max(0, start-chunkStart)
		forwardEnd := min(int64(len(chunkData)), end-chunkStart+1)

		// fmt.Printf("Chunk: %d, chunkStart: %d, forwardStart: %d, forwardEnd: %d\n", i, chunkStart, forwardStart, forwardEnd)

		buf.Write(chunkData[forwardStart:forwardEnd])
	}

	// Stream the chunk data to the response
	if _, err := c.Writer.Write(buf.Bytes()); err != nil {
		c.Error(err)
		return
	}

	c.Writer.Flush()
}
