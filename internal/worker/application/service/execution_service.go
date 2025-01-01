package service

import (
	"context"
	"faas/internal/features/executions/domain/entity"
	"faas/internal/worker/domain/ports"
	"time"
)

type ExecutionService struct {
	containerManager ports.ContainerManager
	executionRepo    ports.ExecutionRepository
}

func NewExecutionService(
	containerManager ports.ContainerManager,
	executionRepo ports.ExecutionRepository,
) *ExecutionService {
	return &ExecutionService{
		containerManager: containerManager,
		executionRepo:    executionRepo,
	}
}

func (s *ExecutionService) ProcessExecution(ctx context.Context, execution *entity.Execution) error {
	// 1. Update status to "running"
	now := time.Now()
	execution.Status = entity.StatusRunning
	execution.StartedAt = &now

	if err := s.executionRepo.UpdateExecution(ctx, execution); err != nil {
		return err
	}

	// 2. Execute function
	output, err := s.containerManager.RunFunction(ctx, execution)
	now = time.Now()
	execution.CompletedAt = &now

	if err != nil {
		// 3a. If there is an error, update status to "failed"
		execution.Status = entity.StatusFailed
		execution.Error = err.Error()
	} else {
		// 3b. If no error, update status to "completed"
		execution.Status = entity.StatusCompleted
		execution.Output = output
	}

	// 4. Save final result
	return s.executionRepo.UpdateExecution(ctx, execution)
}
