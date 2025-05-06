package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"skybox-backend/internal/api/models"
)

type SearchService struct {
	fileRepository   models.FileRepository
	folderRepository models.FolderRepository
	userRepository   models.UserRepository
}

// NewSearchService creates a new instance of SearchService
func NewSearchService(fileRepo models.FileRepository, folderRepo models.FolderRepository, userRepo models.UserRepository) *SearchService {
	return &SearchService{
		fileRepository:   fileRepo,
		folderRepository: folderRepo,
		userRepository:   userRepo,
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

	// Collect all owner IDs from files and folders
	ownerIDs := make(map[string]struct{})
	for _, file := range files {
		ownerIDs[file.OwnerID.Hex()] = struct{}{}
	}
	for _, folder := range folders {
		ownerIDs[folder.OwnerID.Hex()] = struct{}{}
	}

	// Convert ownerIDs map keys to a slice
	ownerIDList := make([]string, 0, len(ownerIDs))
	for id := range ownerIDs {
		ownerIDList = append(ownerIDList, id)
	}

	// Fetch all owners at once
	owners, err := ss.userRepository.GetUsersByIDs(ctx, ownerIDList)
	if err != nil {
		return nil, err
	}

	// Create a map of owner ID to user for quick lookup
	ownerDict := make(map[string]*models.User)
	for _, owner := range owners {
		ownerDict[owner.ID.Hex()] = owner
	}

	// Add owner_email and owner_username to files
	for _, file := range files {
		if owner, exists := ownerDict[file.OwnerID.Hex()]; exists {
			file.OwnerEmail = owner.Email
			file.OwnerUsername = owner.Username
		}
	}

	// Add owner_email and owner_username to folders
	for _, folder := range folders {
		if owner, exists := ownerDict[folder.OwnerID.Hex()]; exists {
			folder.OwnerEmail = owner.Email
			folder.OwnerUsername = owner.Username
		}
	}

	// Combine results
	results := map[string]interface{}{
		"files":   files,
		"folders": folders,
	}

	return results, nil
}
