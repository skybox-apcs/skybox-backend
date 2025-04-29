package controllers

import (
	"bytes"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"skybox-backend/internal/blockserver/services"
	"skybox-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	// Add any dependencies you need here, such as a service or repository
	UploadService *services.UploadService
}

func NewUploadController(uploadService *services.UploadService) *UploadController {
	return &UploadController{
		UploadService: uploadService,
	}
}

const DefaultChunkSize = 5 * 1024 * 1024 // 5MB
const MaxChunkSize = 100 * 1024 * 1024   // 100MB

// UploadWholeFileHandler handles whole file uploads (not chunked)
func (uc *UploadController) UploadWholeFileHandler(c *gin.Context) {
	fileId := c.Param("fileId")
	if fileId == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "Missing file ID")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		shared.ErrorJSON(c, http.StatusBadRequest, "Failed to get file from form")
		return
	}
	defer file.Close()

	// Save the file as-is to S3 or any other storage
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to read file")
		return
	}

	fileName := header.Filename
	ext := filepath.Ext(fileName)

	// Save the file to S3 or any other storage
	err = uc.UploadService.SaveChunk(c, fileId, fileName, ext, 0, buf.Bytes())
	if err != nil {
		shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Return success response
	shared.SuccessJSON(c, http.StatusOK, "File uploaded successfully", gin.H{
		"fileId":   fileId,
		"fileName": fileName,
	})
}

// UploadAutoChunkHandler handles whoel file uploads and split them into chunks
func (uc *UploadController) UploadAutoChunkHandler(c *gin.Context) {
	// Get the fileID from the query parameters
	fileId := c.Param("fileId")
	if fileId == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "Missing file ID")
		return
	}

	// Parse chunk size from query or use default
	chunkSize := DefaultChunkSize
	if size := c.Query("chunkSize"); size != "" {
		if parsedSize, err := strconv.Atoi(size); err == nil && parsedSize > 0 && parsedSize <= MaxChunkSize {
			chunkSize = parsedSize
		}
	}

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		shared.ErrorJSON(c, http.StatusBadRequest, "Failed to get file from form")
	}
	defer file.Close()

	fileName := header.Filename
	ext := filepath.Ext(fileName)

	// Read and split file into chunks
	buf := make([]byte, chunkSize)
	chunkIndex := 0
	for {
		n, err := io.ReadFull(file, buf)
		// If the error is io.EOF, it means we reached the end of the file
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			if n > 0 {
				// Save the last chunk
				err := uc.UploadService.SaveChunk(c, fileId, fileName, ext, chunkIndex, buf[:n])
				if err != nil {
					shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to save chunk")
					return
				}
			}
			break
		}

		if err != nil {
			shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to read file")
			return
		}

		// Save the chunk to S3 or any other storage
		err = uc.UploadService.SaveChunk(c, fileId, fileName, ext, chunkIndex, buf[:n])
		if err != nil {
			shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to save chunk")
			return
		}
		chunkIndex++
	}

	// Return success response
	shared.SuccessJSON(c, http.StatusOK, "File uploaded successfully", gin.H{
		"fileName":   fileName,
		"chunkCount": chunkIndex + 1,
	})
}

// UploadChunkHandler handles chunk uploads
func (uc *UploadController) UploadChunkHandler(c *gin.Context) {
	sessionID := c.Query("sessionID")
	chunkIndexStr := c.Query("chunkIndex")

	// Validate sessionID and chunkIndex
	if sessionID == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "Missing session ID")
		return
	}

	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil {
		shared.ErrorJSON(c, http.StatusBadRequest, "Invalid chunk index")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		shared.ErrorJSON(c, http.StatusBadRequest, "Failed to get file from form")
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to read file")
		return
	}

	// Save the chunk to S3 or any other storage
	err = uc.UploadService.SaveChunkFromSession(c, sessionID, chunkIndex, buf.Bytes())
	if err != nil {
		shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to save chunk")
		return
	}

	// Return success response
	shared.SuccessJSON(c, http.StatusOK, "Chunk uploaded successfully", gin.H{
		"sessionID":  sessionID,
		"chunkIndex": chunkIndexStr,
		"chunkSize":  buf.Len(),
		"fileName":   header.Filename,
	})
}
