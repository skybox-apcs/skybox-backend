package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"skybox-backend/internal/shared"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	// Add any dependencies you need here, such as a service or repository
}

func NewUploadController() *UploadController {
	return &UploadController{
		// Initialize dependencies here
	}
}

const DefaultChunkSize = 5 * 1024 * 1024 // 5MB
const MaxChunkSize = 100 * 1024 * 1024   // 100MB

// UploadWholeFileHandler handles whoel file uploads and split them into chunks
func (uc *UploadController) UploadWholeFileHandler(c *gin.Context) {
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
				if err := saveChunkAuto(fileName, ext, chunkIndex, buf[:n]); err != nil {
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
		if err := saveChunkAuto(fileName, ext, chunkIndex, buf[:n]); err != nil {
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
	if err := saveChunkManual(header.Filename, filepath.Ext(header.Filename), chunkIndex, buf.Bytes()); err != nil {
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

// saveChunkAuto is a helper function to save a chunk of the file
func saveChunkAuto(fileName string, ext string, chunkIndex int, buf []byte) error {
	// TODO: Save to S3 or any other storage
	fmt.Printf("Saving chunk %d of file %s.%s\n", chunkIndex, fileName, ext)
	return nil
}

// saveChunkManual is a helper function to save a chunk of the file manually
func saveChunkManual(fileName string, ext string, chunkIndex int, buf []byte) error {
	// TODO: Save to S3 or any other storage
	fmt.Printf("Saving chunk %d of file %s.%s\n", chunkIndex, fileName, ext)
	return nil
}
