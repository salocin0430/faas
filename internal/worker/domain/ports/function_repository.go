package ports

import (
	"context"
	"faas/internal/features/functions/domain/entity"
)

type FunctionRepository interface {
	GetByID(ctx context.Context, id string) (*entity.Function, error)
}
