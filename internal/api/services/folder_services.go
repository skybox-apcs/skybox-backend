package services

import (
	"context"

	"skybox-backend/internal/api/models"
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

func (fr *FolderService) CreateFolder(ctx context.Context, folder *models.Folder) (*models.Folder, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.CreateFolder(ctx, folder)
}

func (fr *FolderService) GetFolderByID(ctx context.Context, id string) (*models.Folder, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.GetFolderByID(ctx, id)
}

func (fr *FolderService) GetFolderParentIDByFolderID(ctx context.Context, folderID string) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.GetFolderParentIDByFolderID(ctx, folderID)
}

func (fr *FolderService) GetFolderListInFolder(ctx context.Context, id string) ([]*models.Folder, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.GetFolderListInFolder(ctx, id)
}

func (fr *FolderService) GetFileListInFolder(ctx context.Context, folderID string) ([]*models.File, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.GetFileListInFolder(ctx, folderID)
}

func (fr *FolderService) DeleteFolder(ctx context.Context, id string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.DeleteFolder(ctx, id)
}
