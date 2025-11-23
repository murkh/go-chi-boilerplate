package memory

import (
	"context"
	"sync"

	"github.com/murkh/gig-mobile-backend/internal/core/domain"
	"github.com/murkh/gig-mobile-backend/internal/core/ports"
	"github.com/murkh/gig-mobile-backend/internal/platform/errors"
)

type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewUserRepository() ports.UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *UserRepository) Get(ctx context.Context, id string, opts ...ports.QueryOptions) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, errors.NotFound("user not found")
	}
	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.ID]; ok {
		return errors.BadRequest("user already exists")
	}
	r.users[user.ID] = user
	return nil
}

func (r *UserRepository) List(ctx context.Context, opts ...ports.QueryOptions) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}
