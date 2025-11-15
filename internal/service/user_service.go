package service

import (
	"context"

	"github.com/google/uuid"

	"Avito/internal/domain"
	"Avito/internal/repository"
)

type UserServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (s *UserServiceImpl) Create(ctx context.Context, username string) (*domain.User, error) {
	user := &domain.User{
		UserID:   uuid.New().String(),
		Username: username,
		IsActive: true,
	}
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserServiceImpl) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserServiceImpl) SetActive(ctx context.Context, userID string, isActive bool) error {
	return s.userRepo.SetActive(ctx, userID, isActive)
}
