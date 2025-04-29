package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"skybox-backend/configs"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/blockserver/storage"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

// Package services provides the service layer for the controllers
type UploadService struct {
	// Add any dependencies you need here, e.g. repositories, clients, etc.
	baseURL  string
	s3Client *s3.Client
}

func NewUploadService() *UploadService {
	s3client := storage.GetS3Client()
	if s3client == nil && configs.Config.AWSEnabled {
		panic("Failed to create AWS S3 client")
	}

	baseURL := fmt.Sprintf("http://%s:%s", configs.Config.ServerHost, configs.Config.ServerPort)

	return &UploadService{
		// Initialize dependencies here
		baseURL:  baseURL,
		s3Client: s3client,
	}
}

// FetchFileObject is a helper function to retrieve FileObject from API Server
// This helper function would be used in the controller to fetch the file object from the API server
func (us *UploadService) FetchFileObject(ctx *gin.Context, fileId string) (*models.FileResponse, error) {
	// Define the API Server URL
	apiServerURL := fmt.Sprintf("%s/api/v1/files/%s", us.baseURL, fileId)

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiServerURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add authentication header
	headerToken := ctx.GetHeader("Authorization")
	req.Header.Set("Authorization", headerToken)

	// Send the HTTP Request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch file object: %s", resp.Status)
	}

	type responseStruct struct {
		Status  string               `json:"status"`
		Message string               `json:"message"`
		Data    *models.FileResponse `json:"data"`
	}

	response := &responseStruct{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return response.Data, nil
}

// ValidateFile is a helper function to validate the file object
func (us *UploadService) ValidateFile(ctx *gin.Context, fileId string) error {
	// Fetch the file metadata from the API Server
	file, err := us.FetchFileObject(ctx, fileId)
	if err != nil {
		return fmt.Errorf("failed to fetch file object: %w", err)
	}

	// Get the user ID from the context
	userId, ok := ctx.Value("x-user-id").(string)
	if !ok || userId != file.OwnerID {
		return fmt.Errorf("user ID does not match the file owner")
	}

	// Check the status of the file
	if file.Status == "uploaded" {
		return fmt.Errorf("file is already uploaded")
	}

	return nil
}

// SaveChunk is a helper function to save a chunk of the file
// It would upload the chunk to S3 if in production or save it locally if in development
// Key format: `<userId>/<fileId>_<chunkIndex>`
// The fileId is the ID of the file being uploaded, and chunkIndex is the index of the chunk
func (us *UploadService) SaveChunk(ctx context.Context, fileId string, fileName string, ext string, chunkIndex int, buf []byte) error {
	// Get the user id from the ctx which passed from middleware
	// From gin: ctx.GetHeader("x-user-id")
	userId := ctx.Value("x-user-id").(string)
	fmt.Printf("Uploader's user ID: %s\n", userId)
	if userId == "" {
		return fmt.Errorf("missing user ID in context")
	}

	if configs.Config.AWSEnabled {
		// Save to S3
		key := fmt.Sprintf("%s/%s_%d%s", userId, fileId, chunkIndex, ext)
		_, err := us.s3Client.PutObject(
			context.TODO(),
			&s3.PutObjectInput{
				Bucket:      aws.String(configs.Config.AWSBucket),
				Key:         aws.String(key),
				Body:        bytes.NewReader(buf),
				ContentType: aws.String("application/octet-stream"),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to upload chunk to S3: %w", err)
		}

		fmt.Printf("Uploaded chunk %d of file %s to S3 bucket %s\n", chunkIndex, fileId, configs.Config.AWSBucket)
	} else {
		// Save locally (for development/testing purposes)
		// Create the directory if it doesn't exist
		localDir := fmt.Sprintf("tmp/%s", userId)
		err := os.MkdirAll(localDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create local directory: %w", err)
		}

		localPath := fmt.Sprintf("tmp/%s/%s_%d%s", userId, fileId, chunkIndex, ext)
		err = os.WriteFile(localPath, buf, 0644)
		if err != nil {
			return fmt.Errorf("failed to save chunk locally: %w", err)
		}

		fmt.Printf("Saved chunk %d of file %s to %s\n", chunkIndex, fileId, localPath)
	}

	// TODO: Once the chunk is saved, call API Server to update the file status
	fmt.Printf("Calling API Server to update file status...\n")

	return nil
}

// SaveChunkFromSession is a helper function to save a chunk of the file from a session
func (us *UploadService) SaveChunkFromSession(ctx context.Context, sessionId string, chunkIndex int, buf []byte) error {
	// Get the user id from the ctx which passed from middleware
	// From gin: ctx.GetHeader("x-user-id")
	userId := ctx.Value("x-user-id").(string)
	if userId == "" {
		return fmt.Errorf("missing user ID in context")
	}

	if configs.Config.AWSEnabled {
		// Save to S3

	} else {
		// Save locally (for development/testing purposes)

	}

	return nil
}
