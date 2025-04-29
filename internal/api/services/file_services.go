package services

import (
	"context"
	"fmt"
	"math"

	"skybox-backend/configs"
	"skybox-backend/internal/api/models"

	"github.com/google/uuid"
)

// FileService is the service for file operations
type FileService struct {
	fileRepository          models.FileRepository
	uploadSessionRepository models.UploadSessionRepository
}

// NewFileService creates a new instance of the FileService
func NewFileService(fr models.FileRepository, usr models.UploadSessionRepository) *FileService {
	return &FileService{
		fileRepository:          fr,
		uploadSessionRepository: usr,
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
	var maxFileSize = configs.Config.MaxChunkSize
	var uploadURL string = ""

	// If the file is large, create upload session and chunk records
	if true || file.Size > maxFileSize {
		// Create upload session and chunk records
		// uploadURL = fmt.Sprintf("%s:%s/upload/chunked/%s",
		// 	configs.Config.BlockServerHost,
		// 	configs.Config.BlockServerPort,
		// 	savedFile.ID.Hex(),
		// )

		// Create a session for chunked uploads
		totalChunks := int(math.Ceil(float64(file.Size) / float64(configs.Config.DefaultChunkSize)))
		sessionToken := uuid.New().String()
		uploadSession := &models.UploadSession{
			FileID:       savedFile.ID,
			UserID:       savedFile.OwnerID,
			SessionToken: sessionToken,
			TotalSize:    file.Size,
			ChunkSize:    configs.Config.DefaultChunkSize,
			ChunkList:    []int{}, // This will be updated later when chunks are uploaded
			TotalChunks:  totalChunks,
			Status:       "pending",
		}

		_, err = fr.uploadSessionRepository.CreateSessionRecord(ctx, uploadSession)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create upload session: %w", err)
		}

		// Create the uploadURL based on the uploadSession ID
		uploadURL = fmt.Sprintf("%s:%s/upload/session/%s/chunk",
			configs.Config.BlockServerHost,
			configs.Config.BlockServerPort,
			sessionToken,
		)
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
