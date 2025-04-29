package services

import (
	"context"
	"fmt"

	"skybox-backend/configs"
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

// UploadFileMetadata uploads the metadata of a file and returns the saved file and upload URL
// It also decides whether to create an upload session for chunked uploads or a single upload URL for whole file uploads
// TODO: Handle concurrency and chunked uploads
func (fr *FileService) UploadFileMetadata(ctx context.Context, file *models.File) (*models.File, string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	savedFile, err := fr.fileRepository.UploadFileMetadata(ctx, file)
	if err != nil {
		return nil, "", err
	}

	// Decided if chunked or whole file upload
	const maxFileSize = 50 * 1024 * 1024 // 50MB. TODO: Move to config
	var uploadURL string = ""

	// If the file is large, create upload session and chunk records
	if file.Size > maxFileSize {
		// Create upload session and chunk records

	} else {
		// Create a single upload URL for the whole file
		uploadURL = fmt.Sprintf("%s:%s/upload/whole/%s",
			configs.Config.BlockServerHost,
			configs.Config.BlockServerPort,
			savedFile.ID.Hex(),
		)
	}

	return savedFile, uploadURL, nil
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
