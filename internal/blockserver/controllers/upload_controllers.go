package controllers

import (
	"bytes"
	"io"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"

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

const maxWorkers = 5                     // Number of concurrent workers for chunk uploads (TODO: make this configurable)
const DefaultChunkSize = 5 * 1024 * 1024 // 5MB
const MaxChunkSize = 100 * 1024 * 1024   // 100MB

type UploadChunk struct {
	Index int    `json:"index"` // Index of the chunk
	Data  []byte `json:"data"`  // Data of the chunk
}

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

	// Check if the fild is too large
	if header.Size > MaxChunkSize {
		shared.ErrorJSON(c, http.StatusBadRequest, "File size exceeds the maximum limit")
		return
	}

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
		shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to save file"+err.Error())
		return
	}

	// Return success response
	shared.SuccessJSON(c, http.StatusOK, "File uploaded successfully", gin.H{
		"fileId":   fileId,
		"fileName": fileName,
	})
}

// UploadAutoChunkHandler handles whole file uploads and split them into chunks
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

	// Get file from form data
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		shared.ErrorJSON(c, http.StatusBadRequest, "Failed to get file from form")
	}
	defer file.Close()

	fileName := header.Filename
	fileExt := filepath.Ext(fileName)
	totalSize := header.Size
	totalChunks := int(math.Ceil(float64(totalSize) / float64(chunkSize)))

	// Worker pool for concurrent chunk uploads
	uploadChan := make(chan UploadChunk, totalChunks)
	var wg sync.WaitGroup

	for i := 0; i < maxWorkers; i++ {
		go func() {
			for chunk := range uploadChan {
				// Save each chunk to S3 or any other storage
				err := uc.UploadService.SaveChunk(c, fileId, fileName, fileExt, chunk.Index, chunk.Data)
				if err != nil {
					shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to save chunk"+err.Error())
					return
				}
				// Signal that the upload is done
				wg.Done()
			}
		}()
	}

	// Read the file and split it into chunks
	chunkIndex := 0
	for {
		buf := make([]byte, chunkSize)
		n, err := io.ReadFull(file, buf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			if n > 0 {
				wg.Add(1)
				uploadChan <- UploadChunk{
					Index: chunkIndex,
					Data:  buf[:n],
				}
			}
			break
		}

		if err != nil {
			shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to read file"+err.Error())
			return
		}

		// Send chunk to the worker pool
		wg.Add(1)
		uploadChan <- UploadChunk{
			Index: chunkIndex,
			Data:  buf[:n],
		}
		chunkIndex++
	}

	// Close the channel and wait for all uploads to finish
	close(uploadChan)
	wg.Wait()
	if err := file.Close(); err != nil {
		shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to close file"+err.Error())
		return
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
