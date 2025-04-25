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
