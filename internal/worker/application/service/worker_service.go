package service

import (
	"faas/internal/features/executions/domain/entity"
	"faas/internal/worker/domain/processor"
	"faas/internal/worker/domain/repository"
)

type WorkerService struct {
	containerRepo repository.ContainerRepository
	processor     processor.ExecutionProcessor
}

func NewWorkerService(repo repository.ContainerRepository, proc processor.ExecutionProcessor) *WorkerService {
	return &WorkerService{
		containerRepo: repo,
		processor:     proc,
	}
}

func (s *WorkerService) ProcessExecution(execution *entity.Execution) error {
	return s.processor.ProcessExecution(execution)
}
