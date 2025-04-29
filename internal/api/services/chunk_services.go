package services

import (
	"context"

	"skybox-backend/internal/api/models"
)

// ChunkService is the service for chunk management
type ChunkService struct {
	chunkRepository models.ChunkRepository
}

// NewChunkService creates a new instance of the ChunkService
func NewChunkService(cr models.ChunkRepository) *ChunkService {
	return &ChunkService{
		chunkRepository: cr,
	}
}

// UploadChunkMetadata uploads chunk metadata to the database
func (cs *ChunkService) UploadChunkMetadata(ctx context.Context, fileId string, chunk *models.Chunk) (*models.Chunk, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return cs.chunkRepository.UploadChunkMetadata(fileId, chunk)
}

// UpdateChunkStatus updates the status of a chunk in the database
func (cs *ChunkService) UpdateChunkStatus(ctx context.Context, fileId string, chunkIndex int, status string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return cs.chunkRepository.UpdateChunkStatus(fileId, chunkIndex, status)
}

// GetChunksByFileID retrieves all chunks for a specific file from the database
func (cs *ChunkService) GetChunksByFileID(ctx context.Context, fileId string) ([]models.Chunk, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return cs.chunkRepository.GetChunksByFileID(fileId)
}

// GetChunkByID retrieves a chunk by its ID from the database
func (cs *ChunkService) GetChunkByID(ctx context.Context, id string) (*models.Chunk, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return cs.chunkRepository.GetChunkByID(id)
}
