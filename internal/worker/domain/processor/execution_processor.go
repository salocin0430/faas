package processor

import "faas/internal/features/executions/domain/entity"

type ExecutionProcessor interface {
	ProcessExecution(execution *entity.Execution) error
	HandleResult(result *entity.ExecutionResult) error
}
