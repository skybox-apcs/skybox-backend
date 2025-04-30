package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"skybox-backend/internal/api/models"
)

type SearchService struct {
	fileRepository   models.FileRepository
	folderRepository models.FolderRepository
}

// NewSearchService creates a new instance of SearchService
func NewSearchService(fileRepo models.FileRepository, folderRepo models.FolderRepository) *SearchService {
	return &SearchService{
		fileRepository:   fileRepo,
		folderRepository: folderRepo,
	}
}

// SearchFilesAndFolders searches for files and folders based on the query
func (ss *SearchService) SearchFilesAndFolders(ctx context.Context, ownerId primitive.ObjectID, query string) (map[string]interface{}, error) {
	// Search files
	files, err := ss.fileRepository.SearchFiles(ctx, ownerId, query)
	if err != nil {
		return nil, err
	}

	// Search folders
	folders, err := ss.folderRepository.SearchFolders(ctx, ownerId, query)
	if err != nil {
		return nil, err
	}

	// Ensure files and folders are not nil
	if files == nil {
		files = []*models.File{}
	}
	if folders == nil {
		folders = []*models.Folder{}
	}

	// Combine results
	results := map[string]interface{}{
		"files":   files,
		"folders": folders,
	}

	return results, nil
}
