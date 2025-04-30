package services

import (
	"context"
	"fmt"
	"os"

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

	// Create a session for chunked uploads
	sessionToken := uuid.New().String()
	uploadSession := &models.UploadSession{
		FileID:       savedFile.ID,
		UserID:       savedFile.OwnerID,
		SessionToken: sessionToken,
		TotalSize:    file.Size,
		ActualSize:   0,
		ChunkList:    []int{}, // This will be updated later when chunks are uploaded
		Status:       "pending",
	}

	_, err = fr.uploadSessionRepository.CreateSessionRecord(ctx, uploadSession)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create upload session: %w", err)
	}

	// Create the uploadURL based on the uploadSession ID
	uploadURL := fmt.Sprintf("%s:%s/upload/session/%s/chunk",
		configs.Config.BlockServerHost,
		configs.Config.BlockServerPort,
		sessionToken,
	)

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

// GetFileData retrieves the file data for a given file ID and chunk number
// It returns the file data as a byte slice and an error if any occurs
func (fr *FileService) GetFileData(ctx context.Context, fileId string, chunkNumber int) ([]byte, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Get the file metadata
	file, err := fr.fileRepository.GetFileByID(ctx, fileId)
	if err != nil {
		return nil, err
	}

	// Check if the chunk number is valid
	if chunkNumber < 0 || chunkNumber >= file.TotalChunks {
		return nil, fmt.Errorf("invalid chunk number: %d", chunkNumber)
	}

	// Get the file data for the specified chunk number
	key := fmt.Sprintf("%s/%s_%d", file.OwnerID.Hex(), file.ID.Hex(), chunkNumber)
	data := []byte{} // Placeholder for the actual data retrieval

	if configs.Config.AWSEnabled {
		// TODO: Implement AWS S3 chunk retrieval
	} else {
		// Local file retrieval
		// Use local directorty for testing purposes
		full_key := fmt.Sprintf("tmp/%s", key)

		data, err = os.ReadFile(full_key)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve file data: %w", err)
		}
	}

	return data, nil
}
