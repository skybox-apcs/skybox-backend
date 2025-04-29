package controllers

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"skybox-backend/configs"
	"skybox-backend/internal/blockserver/models"
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

var (
	maxWorkers       = configs.Config.MaxWorkers       // Number of concurrent workers for chunk uploads
	DefaultChunkSize = configs.Config.DefaultChunkSize // Default chunk size (5MB)
	MaxChunkSize     = configs.Config.MaxChunkSize     // Maximum chunk size (50MB)
)

type UploadChunk struct {
	Index int    `json:"index"` // Index of the chunk
	Data  []byte `json:"data"`  // Data of the chunk
}

func parseContentRange(header string) (start, end, total int64, err error) {
	if !strings.HasPrefix(header, "bytes ") {
		return 0, 0, 0, fmt.Errorf("invalid Content-Range prefix")
	}
	parts := strings.Split(strings.TrimPrefix(header, "bytes "), "/")
	if len(parts) != 2 {
		return 0, 0, 0, fmt.Errorf("invalid Content-Range format")
	}

	rangePart := strings.Split(parts[0], "-")
	if len(rangePart) != 2 {
		return 0, 0, 0, fmt.Errorf("invalid byte range")
	}

	start, err = strconv.ParseInt(rangePart[0], 10, 64)
	if err != nil {
		return
	}
	end, err = strconv.ParseInt(rangePart[1], 10, 64)
	if err != nil {
		return
	}
	total, err = strconv.ParseInt(parts[1], 10, 64)
	return
}

// UploadWholeFileHandler godoc
//
//	@Summary 		Upload a whole file (without chunking)
//	@Description 	Upload a whole file (without chunking) to the server. This is a simple upload endpoint that does not require chunking. It is useful for smaller files or when chunking is not needed, i.e., lower than 50MB.
//	@Tags 			Upload
//	@Accept 		multipart/form-data
//	@Produce 		json
//	@Param 			fileId 	path string true "File ID" default("fileId") example("fileId")
//	@Param 			file 	formData file true "File" default("file") example("file")
//	@Success 		200 {string} string "File uploaded successfully"
//	@Failure 		400 {string} string "Bad Request" "Invalid file ID or Failed to get file from form or file size exceeds the maximum limit"
//	@Failure 		500 {string} string "Internal Server Error" "Failed to save file"
//	@Router 		/upload/whole/{fileId} [post]
//
// UploadWholeFileHandler handles whole file uploads (not chunked)
func (uc *UploadController) UploadWholeFileHandler(c *gin.Context) {
	fileId := c.Param("fileId")
	if fileId == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "Missing file ID")
		return
	}

	// Validate the file
	err := uc.UploadService.ValidateFile(c, fileId)
	if err != nil {
		shared.ErrorJSON(c, http.StatusBadRequest, "Invalid file ID "+err.Error())
		return
	}

	// Get file from form data
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
	response := &models.NonResumableUploadResponse{
		FileID:     fileId,
		FileName:   fileName,
		FileSize:   header.Size,
		ChunkCount: 1,
	}
	shared.SuccessJSON(c, http.StatusOK, "File uploaded successfully", response)
}

// UploadAutoChunkHandler godoc
//
//	@Summary 		Upload a file in chunks (auto chunking)
//	@Description 	Upload a file in chunks (auto chunking) to the server. This endpoint automatically splits the file into chunks and uploads them concurrently. It is useful for larger files or when chunking is needed, i.e., larger than 50MB.
//	@Tags 			Upload
//	@Accept 		multipart/form-data
//	@Produce 		json
//	@Param 			fileId 	path string true "File ID" default("fileId") example("fileId")
//	@Param 			chunkSize 	query int false "Chunk Size" default(5242880) example(5242880)
//	@Param 			file 	formData file true "File" default("file") example("file")
//	@Success 		200 {string} string "File uploaded successfully"
//	@Failure 		400 {string} string "Bad Request" "Invalid file ID or file size exceeds the maximum limit"
//	@Failure 		500 {string} string "Internal Server Error" "Failed to save file"
//	@Router 		/upload/auto-chunk/{fileId} [post]
//
// UploadAutoChunkHandler handles whole file uploads and split them into chunks
func (uc *UploadController) UploadAutoChunkHandler(c *gin.Context) {
	// Get the fileID from the query parameters
	fileId := c.Param("fileId")
	if fileId == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "Missing file ID")
		return
	}

	// Validate the file
	err := uc.UploadService.ValidateFile(c, fileId)
	if err != nil {
		shared.ErrorJSON(c, http.StatusBadRequest, "Invalid file ID "+err.Error())
		return
	}

	// Parse chunk size from query or use default
	chunkSize := DefaultChunkSize
	if size := c.Query("chunkSize"); size != "" {
		if parsedSize, err := strconv.ParseInt(size, 10, 64); err == nil {
			if parsedSize > 0 && parsedSize <= MaxChunkSize {
				chunkSize = parsedSize
			} else {
				shared.ErrorJSON(c, http.StatusBadRequest, "Invalid chunk size")
				return
			}
		} else {
			shared.ErrorJSON(c, http.StatusBadRequest, "Invalid chunk size")
			return
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

	fmt.Printf("Total size: %d, Chunk size: %d, Total chunks: %d\n", totalSize, chunkSize, totalChunks)

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
	response := &models.NonResumableUploadResponse{
		FileID:     fileId,
		FileName:   fileName,
		FileSize:   totalSize,
		ChunkCount: totalChunks,
	}
	shared.SuccessJSON(c, http.StatusOK, "File uploaded successfully", response)
}

// UploadChunkHandler godoc
//
// UploadChunkHandler handles chunk uploads to a resumable session. Using the sessionId and Content-Range header to identify the chunk and its size.
// The chunk is saved to the server or any other storage.
// It is useful for large files that need to be uploaded in chunks, especially when the upload can be interrupted and resumed later.
func (uc *UploadController) UploadChunkHandler(c *gin.Context) {
	sessionToken := c.Param("sessionToken")
	if sessionToken == "" {
		shared.ErrorJSON(c, http.StatusBadRequest, "Missing session ID")
		return
	}

	// Get the Content-Range header to determine the chunk index and size
	contentRange := c.GetHeader("Content-Range")
	start, end, _, err := parseContentRange(contentRange)
	if err != nil {
		shared.ErrorJSON(c, http.StatusBadRequest, "Invalid Content-Range header "+err.Error())
		return
	}

	chunkIndex := int(start / DefaultChunkSize)
	fileId, err := uc.UploadService.ValidateSession(c, sessionToken, chunkIndex)
	if err != nil {
		shared.ErrorJSON(c, http.StatusBadRequest, "Invalid session ID "+err.Error())
		return
	}

	// Read the chunk data from the request body
	chunkData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to read chunk data "+err.Error())
		return
	}
	defer c.Request.Body.Close()

	fmt.Printf("Chunk index: %d, Start: %d, End: %d, Size: %d, Actual Size: %d\n", chunkIndex, start, end, int(end-start+1), len(chunkData))
	if len(chunkData) != int(end-start+1) {
		shared.ErrorJSON(c, http.StatusBadRequest, "Chunk size mismatch")
		return
	}

	// Save the chunk to S3 or any other storage
	err = uc.UploadService.SaveChunk(c, fileId, "", "", chunkIndex, chunkData)
	if err != nil {
		shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to save chunk "+err.Error())
		return
	}

	// Update the session record in the database
	err = uc.UploadService.UpdateSessionRecord(c, sessionToken, chunkIndex)
	if err != nil {
		shared.ErrorJSON(c, http.StatusInternalServerError, "Failed to update session record "+err.Error())
		return
	}

	// Return success response
	shared.SuccessJSON(c, http.StatusOK, "Chunk uploaded successfully", nil)
}
