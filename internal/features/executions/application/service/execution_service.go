package service

import (
	"context"
	"faas/internal/features/executions/application/dto"
	"faas/internal/features/executions/domain/entity"
	"faas/internal/features/executions/domain/repository"
	functionRepo "faas/internal/features/functions/domain/repository"
	"faas/internal/shared/domain/errors"
	"faas/internal/shared/infrastructure/config"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type ExecutionService struct {
	executionRepo       repository.ExecutionRepository
	executionStreamRepo repository.ExecutionStreamRepository
	functionRepo        functionRepo.FunctionRepository
	config              *config.Config
}

func NewExecutionService(repo repository.ExecutionRepository, streamRepo repository.ExecutionStreamRepository, functionRepo functionRepo.FunctionRepository, config *config.Config) *ExecutionService {
	return &ExecutionService{
		executionRepo:       repo,
		executionStreamRepo: streamRepo,
		functionRepo:        functionRepo,
		config:              config,
	}
}

func (s *ExecutionService) CreateExecution(ctx context.Context, req *dto.CreateExecutionRequest, userID string) (*dto.ExecutionResponse, error) {
	// Check execution limit
	count, err := s.executionRepo.GetActiveExecutionCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	maxExecutions, _ := strconv.Atoi(s.config.MaxConcurrentExecutions)
	if count >= maxExecutions {
		return nil, fmt.Errorf("execution limit exceeded: maximum %d concurrent executions", maxExecutions)
	}

	//Get function and validate if this function has the same userID
	function, err := s.functionRepo.GetByID(ctx, req.FunctionID)
	if err != nil {
		return nil, err
	}

	if function.UserID != userID {
		return nil, errors.NewAppError("unauthorized", "Not authorized to execute this function")
	}

	// Create execution
	execution := &entity.Execution{
		ID:         uuid.New().String(),
		FunctionID: req.FunctionID,
		UserID:     userID,
		Status:     entity.StatusPending,
		Input:      req.Input,
		CreatedAt:  time.Now(),
	}

	if err := s.executionRepo.Save(ctx, execution); err != nil {
		return nil, err
	}

	if err := s.executionStreamRepo.PublishPending(execution); err != nil {
		return nil, err
	}

	return dto.NewExecutionResponse(execution), nil
}

func (s *ExecutionService) GetExecution(ctx context.Context, id string, userID string) (*dto.ExecutionResponse, error) {
	execution, err := s.executionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewAppError("execution_not_found", "Execution not found")
	}

	if execution.UserID != userID {
		return nil, errors.NewAppError("unauthorized", "Not authorized to view this execution")
	}

	return dto.NewExecutionResponse(execution), nil
}

func (s *ExecutionService) ListUserExecutions(ctx context.Context, userID string) ([]*dto.ExecutionResponse, error) {
	executions, err := s.executionRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, errors.NewAppError("list_executions_failed", err.Error())
	}

	responses := make([]*dto.ExecutionResponse, len(executions))
	for i, execution := range executions {
		responses[i] = dto.NewExecutionResponse(execution)
	}

	return responses, nil
}
