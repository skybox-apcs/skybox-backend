package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
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

func requestAPIServer(ctx *gin.Context, method string, url string, body interface{}) (*http.Response, error) {
	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create the request body
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add authentication header
	headerToken := ctx.GetHeader("Authorization")
	req.Header.Set("Authorization", headerToken)

	// Set the content type to JSON
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch file object: %s", resp.Status)
	}

	return resp, nil
}

// FetchFileObject is a helper function to retrieve FileObject from API Server
// This helper function would be used in the controller to fetch the file object from the API server
func (us *UploadService) FetchFileObject(ctx *gin.Context, fileId string) (*models.FileResponse, error) {
	// Define the API Server URL
	apiServerURL := fmt.Sprintf("%s/api/v1/files/%s", us.baseURL, fileId)

	// Create an HTTP request
	resp, err := requestAPIServer(ctx, http.MethodGet, apiServerURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	defer resp.Body.Close()

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

// FetchSessionObject is a helper function to retrieve SessionObject from API Server
// This helper function would be used in the controller to fetch the file object from the API server
func (us *UploadService) FetchSessionObject(ctx *gin.Context, sessionToken string) (*models.UploadSession, error) {
	// Define the API Server URL
	apiServerURL := fmt.Sprintf("%s/api/v1/upload/%s", us.baseURL, sessionToken)

	// Create an HTTP request
	resp, err := requestAPIServer(ctx, http.MethodGet, apiServerURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	defer resp.Body.Close()

	type responseStruct struct {
		Status  string                `json:"status"`
		Message string                `json:"message"`
		Data    *models.UploadSession `json:"data"`
	}

	response := &responseStruct{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return response.Data, nil
}

func (us *UploadService) ValidateSession(ctx *gin.Context, sessionId string, chunkIndex int) (string, error) {
	// Fetch the session metadata from the API Server
	session, err := us.FetchSessionObject(ctx, sessionId)
	if err != nil {
		return "", fmt.Errorf("failed to fetch session object: %w", err)
	}

	// Get the user ID from the context
	userId, ok := ctx.Value("x-user-id").(string)
	if !ok || userId != session.UserID.Hex() {
		return "", fmt.Errorf("user ID does not match the session owner")
	}

	// Check if the chunk in the session.ChunkList or not
	if slices.Contains(session.ChunkList, chunkIndex) {
		return "", fmt.Errorf("chunk %d already exists in the session", chunkIndex)
	}

	return session.FileID.Hex(), nil
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
		key := fmt.Sprintf("%s/%s_%d", userId, fileId, chunkIndex)
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

		localPath := fmt.Sprintf("tmp/%s/%s_%d", userId, fileId, chunkIndex)
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

func (us *UploadService) UpdateSessionRecord(ctx *gin.Context, sessionToken string, chunkIndex int) error {
	// Define the API Server URL
	apiServerURL := fmt.Sprintf("%s/api/v1/upload/%s", us.baseURL, sessionToken)

	// Create an HTTP request
	var requestBody = struct {
		ChunkNumber int `json:"chunk_number"`
	}{
		ChunkNumber: chunkIndex,
	}

	resp, err := requestAPIServer(ctx, http.MethodPut, apiServerURL, requestBody)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	type responseStruct struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	response := &responseStruct{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if response.Status != "success" {
		return fmt.Errorf("failed to update session record: %s", response.Message)
	}

	return nil
}
