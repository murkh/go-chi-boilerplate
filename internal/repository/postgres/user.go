package postgres

import (
	"context"

	"github.com/murkh/gig-mobile-backend/internal/core/domain"
	"github.com/murkh/gig-mobile-backend/internal/core/ports"
	"github.com/murkh/gig-mobile-backend/internal/platform/errors"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Get(ctx context.Context, id string, opts ...ports.QueryOptions) (*domain.User, error) {
	var user domain.User
	query := r.db.WithContext(ctx)

	if len(opts) > 0 {
		for _, include := range opts[0].Includes {
			query = query.Preload(include)
		}
	}

	if err := query.First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("user not found")
		}
		return nil, errors.InternalServerError(err)
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return errors.InternalServerError(err)
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context, opts ...ports.QueryOptions) ([]*domain.User, error) {
	var users []*domain.User
	query := r.db.WithContext(ctx)

	if len(opts) > 0 {
		for _, include := range opts[0].Includes {
			query = query.Preload(include)
		}
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, errors.InternalServerError(err)
	}

	return users, nil
}
