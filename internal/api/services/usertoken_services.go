package services

import (
	"context"

	"skybox-backend/internal/api/models"
)

// UserTokenService is the service for user token management
type UserTokenService struct {
	userTokenRepository models.UserTokenRepository
}

// NewUserTokenService creates a new instance of the UserTokenService
func NewUserTokenService(utr models.UserTokenRepository) *UserTokenService {
	return &UserTokenService{
		userTokenRepository: utr,
	}
}

// CreateUserToken creates a new user token
func (uts *UserTokenService) CreateUserToken(ctx context.Context, userToken *models.UserToken) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return uts.userTokenRepository.CreateUserToken(ctx, userToken)
}

// GetUserTokenByID retrieves a user token by its ID
func (uts *UserTokenService) GetUserTokenByID(ctx context.Context, userTokenID string) (*models.UserToken, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return uts.userTokenRepository.GetUserTokenByID(ctx, userTokenID)
}

// DeleteUserToken deletes a user token
func (uts *UserTokenService) DeleteUserToken(ctx context.Context, token string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return uts.userTokenRepository.DeleteUserToken(ctx, token)
}

// DeleteUserTokensByUserID deletes all user tokens for a specific user
func (uts *UserTokenService) DeleteUserTokensByUserID(ctx context.Context, userID string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return uts.userTokenRepository.DeleteUserTokensByUserID(ctx, userID)
}
