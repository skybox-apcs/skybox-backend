package services

import (
	"context"

	"skybox-backend/internal/api/models"
)

// FileService is the service for file operations
type FileService struct {
	fileRepository models.FileRepository
}

// NewFileService creates a new instance of the FileService
func NewFileService(fr models.FileRepository) *FileService {
	return &FileService{
		fileRepository: fr,
	}
}

func (fr *FileService) UploadFileMetadata(ctx context.Context, file *models.File) (*models.File, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.fileRepository.UploadFileMetadata(ctx, file)
}

func (fr *FileService) GetFileByID(ctx context.Context, id string) (*models.File, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.fileRepository.GetFileByID(ctx, id)
}

func (fr *FileService) DeleteFile(ctx context.Context, id string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.fileRepository.DeleteFile(ctx, id)
}

func (fr *FileService) RenameFile(ctx context.Context, id string, newName string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.fileRepository.RenameFile(ctx, id, newName)
}

func (fr *FileService) MoveFile(ctx context.Context, id string, newParentFolderID string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return fr.fileRepository.MoveFile(ctx, id, newParentFolderID)
}
