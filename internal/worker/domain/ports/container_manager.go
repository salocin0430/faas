package ports

import (
	"context"
	"faas/internal/features/executions/domain/entity"
)

type ContainerManager interface {
	RunFunction(ctx context.Context, execution *entity.Execution) (string, error)
	Stop() error
}
