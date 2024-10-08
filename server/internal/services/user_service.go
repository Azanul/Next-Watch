package services

import (
	"context"

	"github.com/Azanul/Next-Watch/internal/models"
	"github.com/Azanul/Next-Watch/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo repository.UserRepositoryInterface
}

func NewUserService(userRepo repository.UserRepositoryInterface) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	userID := uuid.New()

	return s.userRepo.Create(ctx, &models.User{
		ID:    userID,
		Email: user.Email,
		Name:  user.Name,
		Role:  "USER",
	})
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}
