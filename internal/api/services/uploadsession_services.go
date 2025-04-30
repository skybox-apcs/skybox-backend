package services

import (
	"context"
	"skybox-backend/internal/api/models"
)

type UploadSessionService struct {
	uploadSessionRepository models.UploadSessionRepository
}

func NewUploadSessionService(ur models.UploadSessionRepository) *UploadSessionService {
	return &UploadSessionService{
		uploadSessionRepository: ur,
	}
}

// CreateUploadSession creates a new upload session for a file
func (us *UploadSessionService) CreateSessionRecord(ctx context.Context, session *models.UploadSession) (*models.UploadSession, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return us.uploadSessionRepository.CreateSessionRecord(ctx, session)
}

// GetSessionRecord retrieves an upload session by its session token
func (us *UploadSessionService) GetSessionRecord(ctx context.Context, sessionToken string) (*models.UploadSession, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return us.uploadSessionRepository.GetSessionRecord(ctx, sessionToken)
}

// GetSessionRecordByFileID retrieves an upload session by its file ID
func (us *UploadSessionService) GetSessionRecordByFileID(ctx context.Context, fileID string) (*models.UploadSession, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return us.uploadSessionRepository.GetSessionRecordByFileID(ctx, fileID)
}

// AddChunkSessionRecord adds a chunk to an existing upload session
func (us *UploadSessionService) AddChunkSessionRecord(ctx context.Context, sessionToken string, chunkNumber int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return us.uploadSessionRepository.AddChunkSessionRecord(ctx, sessionToken, chunkNumber)
}

// AddChunkSessionRecordByFileID adds a chunk to an existing upload session by file ID
func (us *UploadSessionService) AddChunkSessionRecordByFileID(ctx context.Context, fileID string, chunkNumber int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return us.uploadSessionRepository.AddChunkSessionRecordByFileID(ctx, fileID, chunkNumber)
}
