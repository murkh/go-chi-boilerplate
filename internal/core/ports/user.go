package ports

import (
	"context"

	"github.com/murkh/gig-mobile-backend/internal/core/domain"
)

type UserRepository interface {
	Get(ctx context.Context, id string, opts ...QueryOptions) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	List(ctx context.Context, opts ...QueryOptions) ([]*domain.User, error)
}

type UserService interface {
	GetUser(ctx context.Context, id string, opts ...QueryOptions) (*domain.User, error)
	CreateUser(ctx context.Context, email, name string) (*domain.User, error)
	ListUsers(ctx context.Context, opts ...QueryOptions) ([]*domain.User, error)
}
