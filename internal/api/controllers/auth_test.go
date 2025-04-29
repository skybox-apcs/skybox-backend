// auth_test.go

package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"skybox-backend/internal/api/models"
	"skybox-backend/internal/api/services"
	"skybox-backend/pkg/utils"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mocking the UserRepository and UserTokenRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUsersByIDs(ctx context.Context, user []string) ([]*models.User, error) {
	args := m.Called(ctx, user)
	if users, ok := args.Get(0).([]*models.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) UpdateUserLastLogin(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockUserTokenRepository struct {
	mock.Mock
}

func (m *MockUserTokenRepository) CreateUserToken(ctx context.Context, token *models.UserToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockUserTokenRepository) FindUserToken(ctx context.Context, token string) (*models.UserToken, error) {
	args := m.Called(ctx, token)
	if userToken, ok := args.Get(0).(*models.UserToken); ok {
		return userToken, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserTokenRepository) DeleteUserToken(ctx context.Context, tokenID string) error {
	args := m.Called(ctx, tokenID)
	return args.Error(0)
}

func (m *MockUserTokenRepository) DeleteUserTokensByUserID(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserTokenRepository) GetUserTokenByUserID(ctx context.Context, userID string) (*[]models.UserToken, error) {
	args := m.Called(ctx, userID)
	if tokens, ok := args.Get(0).(*[]models.UserToken); ok {
		return tokens, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserTokenRepository) GetUserTokenByID(ctx context.Context, tokenID string) (*models.UserToken, error) {
	args := m.Called(ctx, tokenID)
	if token, ok := args.Get(0).(*models.UserToken); ok {
		return token, args.Error(1)
	}
	return nil, args.Error(1)
}

// Setup mock services
func setupMockAuthServices() (*AuthController, *MockUserRepository, *MockUserTokenRepository) {
	mockUserRepo := new(MockUserRepository)
	mockUserTokenRepo := new(MockUserTokenRepository)

	authService := services.NewAuthService(mockUserRepo)
	userTokenService := services.NewUserTokenService(mockUserTokenRepo)

	authController := &AuthController{
		AuthService:      authService,
		UserTokenService: userTokenService,
	}

	return authController, mockUserRepo, mockUserTokenRepo
}

// TestRegisterHandler_EmailAlreadyExists tests registration when email already exists
func TestRegisterHandler_EmailAlreadyExists(t *testing.T) {
	r := gin.Default()
	group := r.Group("/")
	authController, mockUserRepo, _ := setupMockAuthServices()

	authGroup := group.Group("/auth")
	{
		authGroup.POST("/register", authController.RegisterHandler)
	}

	// Mock: email already exists
	mockUserRepo.On("GetUserByEmail", mock.Anything, "testuser@example.com").Return(&models.User{Email: "testuser@example.com"}, nil)

	reqBody := map[string]string{
		"email":    "testuser@example.com",
		"password": "password123",
		"username": "testuser",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUserRepo.AssertExpectations(t)
}

// TestRegisterHandler_UsernameAlreadyExists tests registration when username already exists
func TestRegisterHandler_UsernameAlreadyExists(t *testing.T) {
	r := gin.Default()
	group := r.Group("/")
	authController, mockUserRepo, _ := setupMockAuthServices()

	authGroup := group.Group("/auth")
	{
		authGroup.POST("/register", authController.RegisterHandler)
	}

	// Mock: email does not exist, username exists
	mockUserRepo.On("GetUserByEmail", mock.Anything, "testuser@example.com").Return(nil, errors.New("Error: no document found"))
	mockUserRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(&models.User{Username: "testuser"}, nil)

	reqBody := map[string]string{
		"email":    "testuser@example.com",
		"password": "password123",
		"username": "testuser",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUserRepo.AssertExpectations(t)
}

// TestRegisterHandler_Success tests successful registration
func TestRegisterHandler_Success(t *testing.T) {
	r := gin.Default()
	group := r.Group("/")
	authController, mockUserRepo, _ := setupMockAuthServices()

	authGroup := group.Group("/auth")
	{
		authGroup.POST("/register", authController.RegisterHandler)
	}

	// Mock: email and username do not exist, user creation succeeds
	mockUserRepo.On("GetUserByEmail", mock.Anything, "testuser@example.com").Return(nil, errors.New("Error: no document found"))
	mockUserRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(nil, errors.New("Error: no document found"))
	mockUserRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	reqBody := map[string]string{
		"email":    "testuser@example.com",
		"password": "password123",
		"username": "testuser",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	fmt.Println(rr.Body.String()) // Print the response body for debugging

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockUserRepo.AssertExpectations(t)
}

// TestRegisterHandler_InvalidPayload tests registration with invalid payload
func TestRegisterHandler_InvalidPayload(t *testing.T) {
	r := gin.Default()
	group := r.Group("/")
	authController, _, _ := setupMockAuthServices()

	authGroup := group.Group("/auth")
	{
		authGroup.POST("/register", authController.RegisterHandler)
	}

	// Invalid JSON
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLoginHandler_Success(t *testing.T) {
	r := gin.Default()
	group := r.Group("/")
	authController, mockUserRepo, mockUserTokenRepo := setupMockAuthServices()

	authGroup := group.Group("/auth")
	{
		authGroup.POST("/login", authController.LoginHandler)
	}

	// Mock user data
	mockUser := &models.User{
		ID:           primitive.NewObjectID(),
		Email:        "testuser@example.com",
		PasswordHash: "$2a$12$e8RSN64OYSN5W5jMYSkhaeGJ1OFUR2OG.gvOJZEaT/89Lfvy3KUl6", // bcrypt hash for "password123"
		Username:     "testuser",
		RootFolderID: primitive.NewObjectID(),
	}

	// Mock: GetUserByEmail returns the mock user
	mockUserRepo.On("GetUserByEmail", mock.Anything, "testuser@example.com").Return(mockUser, nil)

	// Mock: UpdateUserLastLogin succeeds
	mockUserRepo.On("UpdateUserLastLogin", mock.Anything, mockUser.ID.Hex()).Return(nil)

	// Mock: CreateUserToken succeeds
	mockUserTokenRepo.On("CreateUserToken", mock.Anything, mock.AnythingOfType("*models.UserToken")).Return(nil)

	// Request body
	reqBody := map[string]string{
		"email":    "testuser@example.com",
		"password": "password123",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	fmt.Println(rr.Body.String()) // Print the response body for debugging

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)
	mockUserRepo.AssertExpectations(t)
	mockUserTokenRepo.AssertExpectations(t)
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	r := gin.Default()
	group := r.Group("/")
	authController, mockUserRepo, _ := setupMockAuthServices()

	authGroup := group.Group("/auth")
	{
		authGroup.POST("/login", authController.LoginHandler)
	}

	// Mock: GetUserByEmail returns an error (user not found)
	mockUserRepo.On("GetUserByEmail", mock.Anything, "testuser@example.com").Return(nil, errors.New("user not found"))

	// Request body
	reqBody := map[string]string{
		"email":    "testuser@example.com",
		"password": "password123",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	mockUserRepo.AssertExpectations(t)
}

func TestLoginHandler_InvalidPayload(t *testing.T) {
	r := gin.Default()
	group := r.Group("/")
	authController, _, _ := setupMockAuthServices()

	authGroup := group.Group("/auth")
	{
		authGroup.POST("/login", authController.LoginHandler)
	}

	// Invalid JSON
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLogoutHandler_Success(t *testing.T) {
	r := gin.Default()
	group := r.Group("/")
	authController, _, mockUserTokenRepo := setupMockAuthServices()

	authGroup := group.Group("/auth")
	{
		authGroup.POST("/logout", authController.LogoutHandler)
	}

	// Mock: DeleteUserToken succeeds
	mockUserTokenRepo.On("DeleteUserToken", mock.Anything, "valid-refresh-token").Return(nil)

	// Mock: GetKeyFromToken returns a valid user ID
	patches := gomonkey.ApplyFunc(utils.GetKeyFromToken, func(key string, requestToken string, secret string) (string, error) {
		return "valid-user-id", nil
	})
	defer patches.Reset()

	// Request body
	reqBody := map[string]string{
		"refresh_token": "valid-refresh-token",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/auth/logout", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)
	mockUserTokenRepo.AssertExpectations(t)
}

func TestLogoutHandler_ExpiredToken(t *testing.T) {
	r := gin.Default()
	group := r.Group("/")
	authController, _, mockUserTokenRepo := setupMockAuthServices()

	authGroup := group.Group("/auth")
	{
		authGroup.POST("/logout", authController.LogoutHandler)
	}

	// Mock: GetKeyFromToken returns an error for expired token
	patches := gomonkey.ApplyFunc(utils.GetKeyFromToken, func(key string, requestToken string, secret string) (string, error) {
		return "valid-user-id", nil
	})
	defer patches.Reset()

	// Mock: DeleteUserToken fails due to expired token
	mockUserTokenRepo.On("DeleteUserToken", mock.Anything, "invalid-refresh-token").Return(errors.New("expired token"))

	// Request body
	reqBody := map[string]string{
		"refresh_token": "invalid-refresh-token",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/auth/logout", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockUserTokenRepo.AssertExpectations(t)
}
