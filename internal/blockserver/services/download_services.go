package services

import (
	"bytes"
	"fmt"
	"os"
	"skybox-backend/configs"
	"skybox-backend/internal/api/models"
	"skybox-backend/internal/blockserver/storage"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

type DownloadService struct {
	// Add any dependencies or configurations needed for the download service
	// For example, a repository to fetch file metadata or a storage client to access files
	uploadService *UploadService
	s3Client      *s3.Client
}

func NewDownloadService() *DownloadService {
	s3client := storage.GetS3Client()
	if s3client == nil && configs.Config.AWSEnabled {
		panic("Failed to create AWS S3 client")
	}

	return &DownloadService{
		// Initialize any dependencies or configurations here
		uploadService: NewUploadService(),
		s3Client:      s3client,
	}
}

func (ds *DownloadService) GetFileMetadata(ctx *gin.Context, fileId string) (*models.FileResponse, error) {
	// Implement the logic to retrieve file metadata by its ID
	// This could involve querying a database or accessing a file storage system

	metadata, err := ds.uploadService.FetchFileObject(ctx, fileId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve file metadata: %w", err)
	}

	return metadata, nil
}

func (ds *DownloadService) DownloadFile(ctx *gin.Context, ownerId string, fileId string, chunkNumber int) ([]byte, error) {
	// Get the file data for the specified chunk number
	key := fmt.Sprintf("%s/%s_%d", ownerId, fileId, chunkNumber)
	data := []byte{} // Placeholder for the actual data retrieval
	err := error(nil)

	if configs.Config.AWSEnabled {
		// TODO: Implement AWS S3 chunk retrieval
		// Use the AWS SDK to retrieve the file data from S3

		output, err := ds.s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &configs.Config.AWSBucket,
			Key:    &key,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve file data from S3: %w", err)
		}
		defer output.Body.Close()

		// Read the object content into memory
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(output.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read chunk body: %w", err)
		}

		return buf.Bytes(), nil
	}

	// Local file rtrieval
	// Use local direectorty for testing purposes
	full_key := fmt.Sprintf("tmp/%s", key)

	data, err = os.ReadFile(full_key)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve file data: %w", err)
	}

	return data, nil
}
