package services

import (
	"context"

	"skybox-backend/internal/api/models"
)

type UserService struct {
	userRepository models.UserRepository
}

func NewUserService(ur models.UserRepository) *UserService {
	return &UserService{
		userRepository: ur,
	}
}

func (us *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return us.userRepository.GetUserByID(ctx, id)
}

func (us *UserService) GetUsersByIDs(ctx context.Context, ids []string) ([]*models.User, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return us.userRepository.GetUsersByIDs(ctx, ids)
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return us.userRepository.GetUserByEmail(ctx, email)
}

func (us *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return us.userRepository.GetUserByUsername(ctx, username)
}

func (us *UserService) GetUsersByEmails(ctx context.Context, emails []string) ([]*models.User, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return us.userRepository.GetUsersByEmails(ctx, emails)
}
