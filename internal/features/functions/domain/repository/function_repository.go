package repository

import (
	"context"

	"faas/internal/features/functions/domain/entity"
)

type FunctionRepository interface {
	Save(ctx context.Context, function *entity.Function) error
	GetByID(ctx context.Context, id string) (*entity.Function, error)
	ListByUserID(ctx context.Context, userID string) ([]*entity.Function, error)
	Delete(ctx context.Context, id string) error
}
