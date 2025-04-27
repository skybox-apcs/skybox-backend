package services

import (
	"context"

	"skybox-backend/internal/api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FolderService is the service for folder operations
type FolderService struct {
	folderRepository models.FolderRepository
}

// NewFolderService creates a new instance of the FolderService
func NewFolderService(fr models.FolderRepository) *FolderService {
	return &FolderService{
		folderRepository: fr,
	}
}

func (fr *FolderService) CreateFolder(ctx context.Context, folder *models.Folder, userID primitive.ObjectID) (*models.Folder, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.CreateFolder(ctx, folder, userID)
}

func (fr *FolderService) GetFolderByID(ctx context.Context, id string, userID primitive.ObjectID) (*models.Folder, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.GetFolderByID(ctx, id, userID)
}

func (fr *FolderService) GetFolderParentIDByFolderID(ctx context.Context, folderID string, userID primitive.ObjectID) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.GetFolderParentIDByFolderID(ctx, folderID, userID)
}

func (fr *FolderService) GetFolderListInFolder(ctx context.Context, id string, userID primitive.ObjectID) ([]*models.Folder, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.GetFolderListInFolder(ctx, id, userID)
}

func (fr *FolderService) GetFileListInFolder(ctx context.Context, folderID string, userID primitive.ObjectID) ([]*models.File, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.GetFileListInFolder(ctx, folderID, userID)
}

func (fr *FolderService) DeleteFolder(ctx context.Context, id string, userID primitive.ObjectID) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.DeleteFolder(ctx, id, userID)
}

func (fr *FolderService) RenameFolder(ctx context.Context, id string, newName string, userID primitive.ObjectID) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.RenameFolder(ctx, id, newName, userID)
}

func (fr *FolderService) MoveFolder(ctx context.Context, id string, newParentID string, userID primitive.ObjectID) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.MoveFolder(ctx, id, newParentID, userID)
}
