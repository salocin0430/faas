package repository

import "faas/internal/features/executions/domain/entity"

type ExecutionStreamRepository interface {
	PublishPending(execution *entity.Execution) error
}
