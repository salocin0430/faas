package service

import (
	"context"
	"time"

	"faas/internal/features/functions/application/dto"
	"faas/internal/features/functions/domain/entity"
	"faas/internal/features/functions/domain/repository"
	"faas/internal/shared/domain/errors"

	"github.com/google/uuid"
)

type FunctionService struct {
	functionRepo repository.FunctionRepository
}

func NewFunctionService(repo repository.FunctionRepository) *FunctionService {
	return &FunctionService{
		functionRepo: repo,
	}
}

func (s *FunctionService) CreateFunction(ctx context.Context, req *dto.CreateFunctionRequest, userID string) (*dto.FunctionResponse, error) {
	function := &entity.Function{
		ID:          uuid.New().String(),
		UserID:      userID,
		Name:        req.Name,
		ImageURL:    req.ImageURL,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if err := s.functionRepo.Save(ctx, function); err != nil {
		return nil, err
	}

	return &dto.FunctionResponse{
		ID:       function.ID,
		Name:     function.Name,
		ImageURL: function.ImageURL,
		UserID:   function.UserID,
	}, nil
}

func (s *FunctionService) GetFunction(ctx context.Context, id string, userID string) (*dto.FunctionResponse, error) {
	function, err := s.functionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewAppError("function_not_found", "Function not found")
	}

	if function.UserID != userID {
		return nil, errors.NewAppError("unauthorized", "Not authorized to access this function")
	}

	return dto.NewFunctionResponse(function), nil
}

func (s *FunctionService) ListUserFunctions(ctx context.Context, userID string) ([]*dto.FunctionResponse, error) {
	functions, err := s.functionRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, errors.NewAppError("list_functions_failed", err.Error())
	}

	responses := make([]*dto.FunctionResponse, len(functions))
	for i, function := range functions {
		responses[i] = dto.NewFunctionResponse(function)
	}

	return responses, nil
}

func (s *FunctionService) DeleteFunction(ctx context.Context, id string, userID string) error {
	function, err := s.functionRepo.GetByID(ctx, id)
	if err != nil {
		return errors.NewAppError("function_not_found", "Function not found")
	}

	if function.UserID != userID {
		return errors.NewAppError("unauthorized", "Not authorized to delete this function")
	}

	if err := s.functionRepo.Delete(ctx, id); err != nil {
		return errors.NewAppError("delete_failed", err.Error())
	}

	return nil
}
