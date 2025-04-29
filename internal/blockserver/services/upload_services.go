package services

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"skybox-backend/configs"
	"skybox-backend/internal/blockserver/storage"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Package services provides the service layer for the controllers
type UploadService struct {
	// Add any dependencies you need here, e.g. repositories, clients, etc.
	s3Client *s3.Client
}

func NewUploadService() *UploadService {
	s3client := storage.GetS3Client()
	if s3client == nil && configs.Config.AWSEnabled {
		panic("Failed to create AWS S3 client")
	}

	return &UploadService{
		// Initialize dependencies here
		s3Client: s3client,
	}
}

// SaveChunk is a helper function to save a chunk of the file
// It would upload the chunk to S3 if in production or save it locally if in development
// Key format: `<userId>/<fileId>_<chunkIndex>`
// The fileId is the ID of the file being uploaded, and chunkIndex is the index of the chunk
func (us *UploadService) SaveChunk(ctx context.Context, fileId string, fileName string, ext string, chunkIndex int, buf []byte) error {
	// Get the user id from the ctx which passed from middleware
	// From gin: ctx.GetHeader("x-user-id")
	userId := ctx.Value("x-user-id").(string)
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
		localPath := fmt.Sprintf("tmp/%s/%s_%d%s", userId, fileId, chunkIndex, ext)
		err := os.WriteFile(localPath, buf, 0644)
		if err != nil {
			return fmt.Errorf("failed to save chunk locally: %w", err)
		}
		fmt.Printf("Saved chunk %d of file %s to %s\n", chunkIndex, fileId, localPath)
	}

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
		key := fmt.Sprintf("%s/%s_%d", userId, sessionId, chunkIndex)
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

		fmt.Printf("Uploaded chunk %d of session %s to S3 bucket %s\n", chunkIndex, sessionId, configs.Config.AWSBucket)
	} else {
		// Save locally (for development/testing purposes)
		localPath := fmt.Sprintf("tmp/%s/%s_%d", userId, sessionId, chunkIndex)
		err := os.WriteFile(localPath, buf, 0644)
		if err != nil {
			return fmt.Errorf("failed to save chunk locally: %w", err)
		}
		fmt.Printf("Saved chunk %d of session %s to %s\n", chunkIndex, sessionId, localPath)
	}

	return nil
}
