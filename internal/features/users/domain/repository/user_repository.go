package repository

import (
	"context"
	"faas/internal/features/users/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*entity.User, error)
}
