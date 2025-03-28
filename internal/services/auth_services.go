package services

import (
	"context"

	"skybox-backend/internal/models"
	"skybox-backend/pkg/utils"
)

// AuthService is the service for the authentication
type AuthService struct {
	userRepository models.UserRepository
}

// NewAuthService creates a new instance of the AuthService
func NewAuthService(ur models.UserRepository) *AuthService {
	return &AuthService{
		userRepository: ur,
	}
}

// RegisterUser registers a new user
func (as *AuthService) RegisterUser(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return as.userRepository.CreateUser(ctx, user)
}

// GetUserByEmail retrieves a user by email
func (as *AuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return as.userRepository.GetUserByEmail(ctx, email)
}

// GetUserByID retrieves a user by ID
func (as *AuthService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return as.userRepository.GetUserByID(ctx, id)
}

// CreateAccessToken creates an access token for the user
func (as *AuthService) CreateAccessToken(user *models.User, secret string, expiry int) (string, error) {
	return utils.CreateAccessToken(user, secret, expiry)
}
