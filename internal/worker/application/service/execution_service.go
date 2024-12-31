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
	// 1. Actualizar estado a "running"
	now := time.Now()
	execution.Status = entity.StatusRunning
	execution.StartedAt = &now

	if err := s.executionRepo.UpdateExecution(ctx, execution); err != nil {
		return err
	}

	// 2. Ejecutar función
	output, err := s.containerManager.RunFunction(ctx, execution)
	now = time.Now()
	execution.CompletedAt = &now

	if err != nil {
		// 3a. Si hay error, actualizar estado a "failed"
		execution.Status = entity.StatusFailed
		execution.Error = err.Error()
	} else {
		// 3b. Si no hay error, actualizar estado a "completed"
		execution.Status = entity.StatusCompleted
		execution.Output = output
	}

	// 4. Guardar resultado final
	return s.executionRepo.UpdateExecution(ctx, execution)
}
