package repository

import (
	"context"

	"faas/internal/features/executions/domain/entity"
)

type ExecutionRepository interface {
	Save(ctx context.Context, execution *entity.Execution) error
	GetByID(ctx context.Context, id string) (*entity.Execution, error)
	ListByUserID(ctx context.Context, userID string) ([]*entity.Execution, error)
	Update(ctx context.Context, execution *entity.Execution) error
}
