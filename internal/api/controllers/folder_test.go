package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/services"

	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFolderRepository mocks the FolderService for testing
type MockFolderRepository struct {
	mock.Mock
}

func (m *MockFolderRepository) CreateFolder(ctx context.Context, folder *models.Folder) (*models.Folder, error) {
	args := m.Called(ctx, folder)
	if folder, ok := args.Get(0).(*models.Folder); ok {
		return folder, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFolderRepository) GetFolderByID(ctx context.Context, id string) (*models.Folder, error) {
	args := m.Called(ctx, id)
	if folder, ok := args.Get(0).(*models.Folder); ok {
		return folder, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFolderRepository) GetFolderParentIDByFolderID(ctx context.Context, folderID string) (string, error) {
	args := m.Called(ctx, folderID)
	return args.String(0), args.Error(1)
}

func (m *MockFolderRepository) GetFolderListInFolder(ctx context.Context, id string) ([]*models.Folder, error) {
	args := m.Called(ctx, id)
	if folders, ok := args.Get(0).([]*models.Folder); ok {
		return folders, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFolderRepository) GetFolderResponseListInFolder(ctx context.Context, folderID string) ([]*models.FolderResponse, error) {
	args := m.Called(ctx, folderID)
	if folderResponses, ok := args.Get(0).([]*models.FolderResponse); ok {
		return folderResponses, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFolderRepository) GetFileListInFolder(ctx context.Context, folderID string) ([]*models.File, error) {
	args := m.Called(ctx, folderID)
	if files, ok := args.Get(0).([]*models.File); ok {
		return files, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFolderRepository) GetFileResponseListInFolder(ctx context.Context, folderID string) ([]*models.FileResponse, error) {
	args := m.Called(ctx, folderID)
	if fileResponse, ok := args.Get(0).([]*models.FileResponse); ok {
		return fileResponse, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFolderRepository) DeleteFolder(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFolderRepository) RenameFolder(ctx context.Context, id string, newName string) error {
	args := m.Called(ctx, id, newName)
	return args.Error(0)
}

func (m *MockFolderRepository) MoveFolder(ctx context.Context, id string, newParentID string) error {
	args := m.Called(ctx, id, newParentID)
	return args.Error(0)
}

// MockFileRepository mocks the FileService for testing
type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) UploadFileMetadata(ctx context.Context, file *models.File) (*models.File, error) {
	args := m.Called(ctx, file)
	if file, ok := args.Get(0).(*models.File); ok {
		return file, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFileRepository) GetFileByID(ctx context.Context, id string) (*models.File, error) {
	args := m.Called(ctx, id)
	if file, ok := args.Get(0).(*models.File); ok {
		return file, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFileRepository) DeleteFile(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFileRepository) RenameFile(ctx context.Context, id string, newName string) error {
	args := m.Called(ctx, id, newName)
	return args.Error(0)
}

func (m *MockFileRepository) MoveFile(ctx context.Context, id string, newParentFolderID string) error {
	args := m.Called(ctx, id, newParentFolderID)
	return args.Error(0)
}

// Setup Mock Services
func setupMockServices() (*FolderController, *MockFolderRepository, *MockFileRepository) {
	mockFolderRepo := new(MockFolderRepository)
	mockFileRepo := new(MockFileRepository)

	folderService := services.NewFolderService(mockFolderRepo)
	fileService := services.NewFileService(mockFileRepo)

	folderController := &FolderController{
		FolderService: folderService,
		FileService:   fileService,
	}

	return folderController, mockFolderRepo, mockFileRepo
}

func TestFetchRootFolderContents(t *testing.T) {
	r := gin.Default()
	group := r.Group("/")
	folderController, mockFolderRepo, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.GET("/:folderId/contents", folderController.GetContentsHandler)
	}

	// Mock login process to get bearer token and root_folder_id
	mockRootFolderID := "root_folder_id_123"
	mockBearerToken := "Bearer mock_token"

	// Mock folder contents
	mockFolders := []*models.FolderResponse{
		{
			ID:   "0123456789abcdef12345678",
			Name: "Folder 1",
		},
		{
			ID:   "abcdef012345678912345678",
			Name: "Folder 2",
		},
	}
	mockFiles := []*models.FileResponse{
		{
			ID:   "0123456789abcdef12345679",
			Name: "File 1",
		},
		{
			ID:   "abcdef012345678912345679",
			Name: "File 2",
		},
	}

	// Mock GetFolderListInFolder and GetFileListInFolder
	mockFolderRepo.On("GetFolderResponseListInFolder", mock.Anything, mockRootFolderID).Return(mockFolders, nil)
	mockFolderRepo.On("GetFileResponseListInFolder", mock.Anything, mockRootFolderID).Return(mockFiles, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/folders/"+mockRootFolderID+"/contents", nil)
	req.Header.Set("Authorization", mockBearerToken)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response body
	var response struct {
		Data struct {
			Folders []*models.FolderResponse `json:"folder_list"`
			Files   []*models.FileResponse   `json:"file_list"`
		} `json:"data"`
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert folders and files
	assert.Equal(t, mockFolders, response.Data.Folders)
	assert.Equal(t, mockFiles, response.Data.Files)
	assert.Equal(t, "success", response.Status)

	// Assert mock expectations
	mockFolderRepo.AssertExpectations(t)
}

func TestFetchOtherRootContents(t *testing.T) {
	r := gin.Default()
	group := r.Group("/")
	folderController, mockFolderRepo, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.GET("/:folderId/contents", folderController.GetContentsHandler)
	}

	// Mock login process to get bearer token and root_folder_id
	mockBearerToken := "Bearer mock_token"
	mockOtherRootFolderID := "root_folder_id_456" // This is the root ID of other user

	// Mock GetFolderListInFolder and GetFileListInFolder
	mockFolderRepo.On("GetFolderResponseListInFolder", mock.Anything, mockOtherRootFolderID).Return(nil, fmt.Errorf("folder not found"))
	mockFolderRepo.On("GetFileResponseListInFolder", mock.Anything, mockOtherRootFolderID).Return(nil, fmt.Errorf("folder not found"))
	// Create request
	req, _ := http.NewRequest("GET", "/folders/"+mockOtherRootFolderID+"/contents", nil)
	req.Header.Set("Authorization", mockBearerToken)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Parse response body
	var response struct {
		Data struct {
			Folders []*models.FolderResponse `json:"folder_list"`
			Files   []*models.FileResponse   `json:"file_list"`
		} `json:"data"`
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Assert mock expectations
	mockFolderRepo.AssertExpectations(t)
}
