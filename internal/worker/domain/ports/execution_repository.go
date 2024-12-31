package ports

import (
	"context"
	"faas/internal/features/executions/domain/entity"
)

type ExecutionRepository interface {
	UpdateExecution(ctx context.Context, execution *entity.Execution) error
}
