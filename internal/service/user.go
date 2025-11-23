package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/murkh/gig-mobile-backend/internal/core/domain"
	"github.com/murkh/gig-mobile-backend/internal/core/ports"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetUser(ctx context.Context, id string, opts ...ports.QueryOptions) (*domain.User, error) {
	return s.repo.Get(ctx, id, opts...)
}

func (s *UserService) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
	user := &domain.User{
		ID:        uuid.New().String(),
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context, opts ...ports.QueryOptions) ([]*domain.User, error) {
	return s.repo.List(ctx, opts...)
}
