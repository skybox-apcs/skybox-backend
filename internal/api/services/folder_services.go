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

func (fr *FolderService) GetFolderResponseListInFolder(ctx context.Context, folderID string) ([]*models.FolderResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.GetFolderResponseListInFolder(ctx, folderID)
}

func (fr *FolderService) GetFileListInFolder(ctx context.Context, folderID string) ([]*models.File, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.GetFileListInFolder(ctx, folderID)
}

func (fr *FolderService) GetFileResponseListInFolder(ctx context.Context, folderID string) ([]*models.FileResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.GetFileResponseListInFolder(ctx, folderID)
}

func (fr *FolderService) DeleteFolder(ctx context.Context, id string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.DeleteFolder(ctx, id)
}

func (fr *FolderService) RenameFolder(ctx context.Context, id string, newName string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.RenameFolder(ctx, id, newName)
}

func (fr *FolderService) MoveFolder(ctx context.Context, id string, newParentID string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.folderRepository.MoveFolder(ctx, id, newParentID)
}

func (fs *FolderService) GetFolderSharedUsers(ctx context.Context, folderID string) ([]*models.FolderSharedUser, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fs.folderRepository.GetFolderSharedUsers(ctx, folderID)
}

func (fs *FolderService) GetFolderShareInfo(ctx context.Context, folderID string) (bool, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fs.folderRepository.GetFolderShareInfo(ctx, folderID)
}

func (fs *FolderService) GetFolderSharedUser(ctx context.Context, folderID string, userID string) (*models.FolderSharedUser, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fs.folderRepository.GetFolderSharedUser(ctx, folderID, userID)
}

func (fs *FolderService) UpdateFolderPublicStatus(ctx context.Context, folderID string, isPublic bool) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fs.folderRepository.UpdateFolderPublicStatus(ctx, folderID, isPublic)
}

func (fs *FolderService) UpdateFolderAndSubfoldersPublicStatus(ctx context.Context, folderID string, isPublic bool) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fs.folderRepository.UpdateFolderAndAllSubfoldersPublicStatus(ctx, folderID, isPublic)
}

func (fs *FolderService) ShareFolder(ctx context.Context, folderID, userID string, permission bool) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fs.folderRepository.ShareFolder(ctx, folderID, userID, permission)
}

func (fs *FolderService) RemoveFolderShare(ctx context.Context, folderID, userID string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fs.folderRepository.RemoveFolderShare(ctx, folderID, userID)
}

func (fs *FolderService) ShareFolderAndSubfolders(ctx context.Context, folderID, userID string, permission bool) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fs.folderRepository.ShareFolderAndAllSubfolders(ctx, folderID, userID, permission)
}

func (fs *FolderService) RevokeFolderAndSubfoldersShare(ctx context.Context, folderID, userID string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fs.folderRepository.RevokeFolderAndAllSubfoldersShare(ctx, folderID, userID)
}
