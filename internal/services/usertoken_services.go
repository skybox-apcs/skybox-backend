package services

import (
	"context"

	"skybox-backend/internal/models"
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

// DeleteUserToken deletes a user token
func (uts *UserTokenService) DeleteUserToken(ctx context.Context, userTokenID string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return uts.userTokenRepository.DeleteUserToken(ctx, userTokenID)
}
