package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/services"
	"skybox-backend/internal/shared/middlewares"

	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type MockUploadSessionRepository struct {
	mock.Mock
}

func (m *MockUploadSessionRepository) CreateSessionRecord(ctx context.Context, sessionData *models.UploadSession) (*models.UploadSession, error) {
	args := m.Called(ctx, sessionData)
	if us, ok := args.Get(0).(*models.UploadSession); ok {
		return us, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUploadSessionRepository) GetSessionRecord(ctx context.Context, sessionToken string) (*models.UploadSession, error) {
	args := m.Called(ctx, sessionToken)
	if us, ok := args.Get(0).(*models.UploadSession); ok {
		return us, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUploadSessionRepository) GetSessionRecordByFileID(ctx context.Context, fileId string) (*models.UploadSession, error) {
	args := m.Called(ctx, fileId)
	if us, ok := args.Get(0).(*models.UploadSession); ok {
		return us, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUploadSessionRepository) AddChunkSessionRecord(ctx context.Context, sessionToken string, chunkNumber int, chunkSize int, chunkHash string) error {
	args := m.Called(ctx, sessionToken, chunkNumber, chunkSize, chunkHash)
	return args.Error(0)
}

func (m *MockUploadSessionRepository) AddChunkSessionRecordByFileID(ctx context.Context, fileId string, chunkNumber int, chunkSize int, chunkHash string) error {
	args := m.Called(ctx, fileId, chunkNumber, chunkSize, chunkHash)
	return args.Error(0)
}

// Setup Mock Services
func setupMockServices() (*FolderController, *MockFolderRepository, *MockFileRepository, *MockUploadSessionRepository) {
	mockFolderRepo := new(MockFolderRepository)
	mockFileRepo := new(MockFileRepository)
	mockUSRepository := new(MockUploadSessionRepository)

	folderService := services.NewFolderService(mockFolderRepo)
	fileService := services.NewFileService(mockFileRepo, mockUSRepository)

	folderController := &FolderController{
		FolderService: folderService,
		FileService:   fileService,
	}

	return folderController, mockFolderRepo, mockFileRepo, mockUSRepository
}

func TestFetchRootFolderContents(t *testing.T) {
	r := gin.Default()
	r.Use(middlewares.GlobalErrorMiddleware()) // Use the global error middleware
	group := r.Group("/")
	folderController, mockFolderRepo, _, _ := setupMockServices()

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
	r.Use(middlewares.GlobalErrorMiddleware()) // Use the global error middleware
	group := r.Group("/")
	folderController, mockFolderRepo, _, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.GET("/:folderId/contents", folderController.GetContentsHandler)
	}

	// Mock login process to get bearer token and root_folder_id
	mockBearerToken := "Bearer mock_token"
	mockOtherRootFolderID := "root_folder_id_456" // This is the root ID of other user

	// Mock GetFolderListInFolder and GetFileListInFolder
	mockFolderRepo.On("GetFolderResponseListInFolder", mock.Anything, mockOtherRootFolderID).Return(nil, fmt.Errorf("folder not found or deleted"))
	mockFolderRepo.On("GetFileResponseListInFolder", mock.Anything, mockOtherRootFolderID).Return(nil, fmt.Errorf("folder not found or deleted"))

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

func TestRenameFolder_Success(t *testing.T) {
	r := gin.Default()
	r.Use(middlewares.GlobalErrorMiddleware()) // Use the global error middleware
	group := r.Group("/")
	folderController, mockFolderRepo, _, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.PUT("/:folderId/rename", folderController.RenameFolderHandler)
	}

	// Mock folder ID and new name
	mockFolderID := "folder_id_123"
	mockNewName := "New Folder Name"

	// Mock RenameFolder to succeed
	mockFolderRepo.On("RenameFolder", mock.Anything, mockFolderID, mockNewName).Return(nil)

	// Create request
	reqBody := map[string]string{
		"new_name": mockNewName,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/folders/"+mockFolderID+"/rename", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response body
	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Folder renamed successfully.", response.Message)

	// Assert mock expectations
	mockFolderRepo.AssertExpectations(t)
}

func TestRenameFolder_RootFolderError(t *testing.T) {
	r := gin.Default()
	r.Use(middlewares.GlobalErrorMiddleware()) // Use the global error middleware
	group := r.Group("/")
	folderController, mockFolderRepo, _, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.PUT("/:folderId/rename", folderController.RenameFolderHandler)
	}

	// Mock root folder ID and new name
	mockRootFolderID := "root_folder_id_123"
	mockNewName := "New Root Folder Name"

	// Mock RenameFolder to return an error for root folder
	mockFolderRepo.On("RenameFolder", mock.Anything, mockRootFolderID, mockNewName).Return(fmt.Errorf("cannot rename root folder"))

	// Create request
	reqBody := map[string]string{
		"new_name": mockNewName,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/folders/"+mockRootFolderID+"/rename", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusForbidden, rr.Code)

	// Parse response body
	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)

	// Assert mock expectations
	mockFolderRepo.AssertExpectations(t)
}

func TestRenameFolder_InvalidPayload(t *testing.T) {
	r := gin.Default()
	r.Use(middlewares.GlobalErrorMiddleware()) // Use the global error middleware
	group := r.Group("/")
	folderController, _, _, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.PUT("/:folderId/rename", folderController.RenameFolderHandler)
	}

	// Mock folder ID
	mockFolderID := "folder_id_123"

	// Create request with invalid payload
	req, _ := http.NewRequest("PUT", "/folders/"+mockFolderID+"/rename", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse response body
	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Invalid request.", response.Message)
}

func TestRenameFolder_OtherUserFolder(t *testing.T) {
	r := gin.Default()
	r.Use(middlewares.GlobalErrorMiddleware()) // Use the global error middleware
	group := r.Group("/")
	folderController, mockFolderRepo, _, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.PUT("/:folderId/rename", folderController.RenameFolderHandler)
	}

	// Mock folder ID and new name
	mockFolderID := "folder_id_123"
	mockNewName := "New Folder Name"

	// Mock RenameFolder to return an error for other user's folder
	mockFolderRepo.On("RenameFolder", mock.Anything, mockFolderID, mockNewName).Return(fmt.Errorf("folder not found"))

	// Create request
	reqBody := map[string]string{
		"new_name": mockNewName,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/folders/"+mockFolderID+"/rename", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusNotFound, rr.Code)

}

func TestDeleteFolder_RootFolder(t *testing.T) {
	r := gin.Default()
	r.Use(middlewares.GlobalErrorMiddleware()) // Use the global error middleware
	group := r.Group("/")
	folderController, mockFolderRepo, _, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.DELETE("/:folderId", folderController.DeleteFolderHandler)
	}

	// Mock root folder ID
	mockRootFolderID := "root_folder_id_123"

	// Mock DeleteFolder to return an error for root folder
	mockFolderRepo.On("DeleteFolder", mock.Anything, mockRootFolderID).Return(fmt.Errorf("cannot delete root folder"))

	// Create request
	req, _ := http.NewRequest("DELETE", "/folders/"+mockRootFolderID, nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusForbidden, rr.Code)

	// Parse response body
	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)

	// Assert mock expectations
	mockFolderRepo.AssertExpectations(t)
}

func TestDeleteFolder_OtherUserFolder(t *testing.T) {
	r := gin.Default()
	r.Use(middlewares.GlobalErrorMiddleware()) // Use the global error middleware
	group := r.Group("/")
	folderController, mockFolderRepo, _, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.DELETE("/:folderId", folderController.DeleteFolderHandler)
	}

	// Mock folder ID
	mockFolderID := "folder_id_123"

	// Mock DeleteFolder to return an error for other user's folder
	mockFolderRepo.On("DeleteFolder", mock.Anything, mockFolderID).Return(fmt.Errorf("folder not found"))

	// Create request
	req, _ := http.NewRequest("DELETE", "/folders/"+mockFolderID, nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Parse response body
	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)

	// Assert mock expectations
	mockFolderRepo.AssertExpectations(t)
}

func TestDeleteFolder_Success(t *testing.T) {
	r := gin.Default()
	r.Use(middlewares.GlobalErrorMiddleware()) // Use the global error middleware
	group := r.Group("/")
	folderController, mockFolderRepo, _, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.DELETE("/:folderId", folderController.DeleteFolderHandler)
	}

	// Mock folder ID
	mockFolderID := "folder_id_123"

	// Mock DeleteFolder to succeed
	mockFolderRepo.On("DeleteFolder", mock.Anything, mockFolderID).Return(nil)

	// Create request
	req, _ := http.NewRequest("DELETE", "/folders/"+mockFolderID, nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response body
	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)

	// Assert mock expectations
	mockFolderRepo.AssertExpectations(t)
}

func TestMoveFolder_Success(t *testing.T) {
	r := gin.Default()
	r.Use(middlewares.GlobalErrorMiddleware()) // Use the global error middleware
	group := r.Group("/")
	folderController, mockFolderRepo, _, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.PUT("/:folderId/move", folderController.MoveFolderHandler)
	}

	// Mock folder ID and new parent ID
	mockFolderID := "folder_id_123"
	mockNewParentID := "new_parent_id_456"

	// Mock MoveFolder to succeed
	mockFolderRepo.On("MoveFolder", mock.Anything, mockFolderID, mockNewParentID).Return(nil)

	// Create request
	reqBody := map[string]string{
		"new_parent_id": mockNewParentID,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/folders/"+mockFolderID+"/move", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response body
	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Folder moved successfully.", response.Message)

	// Assert mock expectations
	mockFolderRepo.AssertExpectations(t)
}

func TestMoveFolder_RootFolder(t *testing.T) {
	r := gin.Default()
	r.Use(middlewares.GlobalErrorMiddleware()) // Use the global error middleware
	group := r.Group("/")
	folderController, mockFolderRepo, _, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.PUT("/:folderId/move", folderController.MoveFolderHandler)
	}

	// Mock root folder ID and new parent ID
	mockRootFolderID := "root_folder_id_123"
	mockNewParentID := "new_parent_id_456"

	// Mock MoveFolder to return an error for root folder
	mockFolderRepo.On("MoveFolder", mock.Anything, mockRootFolderID, mockNewParentID).Return(fmt.Errorf("cannot move root folder"))

	// Create request
	reqBody := map[string]string{
		"new_parent_id": mockNewParentID,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/folders/"+mockRootFolderID+"/move", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusForbidden, rr.Code)

	// Parse response body
	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)

	// Assert mock expectations
	mockFolderRepo.AssertExpectations(t)
}

func TestMoveFolder_OtherUserFolder(t *testing.T) {
	r := gin.Default()
	r.Use(middlewares.GlobalErrorMiddleware()) // Use the global error middleware
	group := r.Group("/")
	folderController, mockFolderRepo, _, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.PUT("/:folderId/move", folderController.MoveFolderHandler)
	}

	// Mock folder ID and new parent ID
	mockFolderID := "folder_id_123"
	mockNewParentID := "new_parent_id_456"

	// Mock MoveFolder to return an error for other user's folder
	mockFolderRepo.On("MoveFolder", mock.Anything, mockFolderID, mockNewParentID).Return(fmt.Errorf("folder not found"))

	// Create request
	reqBody := map[string]string{
		"new_parent_id": mockNewParentID,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/folders/"+mockFolderID+"/move", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Parse response body
	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)

	// Assert mock expectations
	mockFolderRepo.AssertExpectations(t)
}

func TestMoveFolder_InvalidPayload(t *testing.T) {
	r := gin.Default()
	r.Use(middlewares.GlobalErrorMiddleware()) // Use the global error middleware
	group := r.Group("/")
	folderController, _, _, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.PUT("/:folderId/move", folderController.MoveFolderHandler)
	}

	// Mock folder ID
	mockFolderID := "folder_id_123"

	// Create request with invalid payload
	req, _ := http.NewRequest("PUT", "/folders/"+mockFolderID+"/move", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse response body
	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Invalid request.", response.Message)
}

func TestCreateFolder_Success(t *testing.T) {
	// Mock folder ID and request body
	mockFolderID := primitive.NewObjectID().Hex()
	mockFolderName := "New Folder"
	mockOwnerID := primitive.NewObjectID()
	mockParentFolderID := primitive.NewObjectID()

	// Mock folder creation request
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("x-user-id-hex", mockOwnerID)
		c.Next()
	})
	r.Use(middlewares.GlobalErrorMiddleware()) // Use the global error middleware
	group := r.Group("/")
	folderController, mockFolderRepo, _, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.POST("/:folderId/create", folderController.CreateFolderHandler)
	}

	// Mock folder creation response
	mockFolder := &models.Folder{
		ID:             primitive.NewObjectID(),
		Name:           mockFolderName,
		OwnerID:        mockOwnerID,
		ParentFolderID: mockParentFolderID,
		IsDeleted:      false,
		IsRoot:         false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Mock CreateFolder to succeed
	mockFolderRepo.On("CreateFolder", mock.Anything, mock.AnythingOfType("*models.Folder")).Return(mockFolder, nil)

	// Create request
	reqBody := map[string]string{
		"name": mockFolderName,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/folders/"+mockFolderID+"/create", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Parse response body
	var response struct {
		Data    models.CreateFolderResponse `json:"data"`
		Message string                      `json:"message"`
		Status  string                      `json:"status"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Folder created successfully.", response.Message)
	assert.Equal(t, mockFolder.ID.Hex(), response.Data.ID)
	assert.Equal(t, mockFolder.OwnerID.Hex(), response.Data.OwnerID)
	assert.Equal(t, mockFolderID, response.Data.ParentFolderID)
	assert.Equal(t, mockFolder.Name, response.Data.Name)

	// Assert mock expectations
	mockFolderRepo.AssertExpectations(t)
}

func TestCreateFolder_InvalidPayload(t *testing.T) {
	r := gin.Default()
	r.Use(middlewares.GlobalErrorMiddleware()) // Use the global error middleware
	group := r.Group("/")
	folderController, _, _, _ := setupMockServices()

	folderGroup := group.Group("/folders")
	{
		folderGroup.POST("/:folderId/create", folderController.CreateFolderHandler)
	}

	// Mock folder ID
	mockFolderID := primitive.NewObjectID().Hex()

	// Create request with invalid payload
	req, _ := http.NewRequest("POST", "/folders/"+mockFolderID+"/create", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Parse response body
	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)
	assert.Equal(t, "Invalid request.", response.Message)
}
