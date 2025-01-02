package service

import (
	"context"
	"faas/internal/features/function_objects/application/dto"
	"faas/internal/features/function_objects/domain/entity"
	"faas/internal/features/function_objects/domain/repository"
	"time"

	"github.com/google/uuid"
)

type ObjectService struct {
	objectRepo repository.ObjectRepository
}

func NewObjectService(repo repository.ObjectRepository) *ObjectService {
	return &ObjectService{
		objectRepo: repo,
	}
}

func (s *ObjectService) CreateObject(ctx context.Context, req *dto.CreateObjectRequest, data []byte, contentType string) (*dto.ObjectResponse, error) {
	obj := &entity.FunctionObject{
		ID:          uuid.New().String(),
		FunctionID:  req.FunctionID,
		Name:        req.Name,
		Size:        int64(len(data)),
		ContentType: contentType,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.objectRepo.Save(ctx, obj, data); err != nil {
		return nil, err
	}

	return dto.NewObjectResponse(obj), nil
}

func (s *ObjectService) GetObject(ctx context.Context, functionID, name string) (*entity.FunctionObject, []byte, error) {
	return s.objectRepo.Get(ctx, functionID, name)
}

func (s *ObjectService) ListObjects(ctx context.Context, functionID string) ([]*dto.ObjectResponse, error) {
	objects, err := s.objectRepo.List(ctx, functionID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ObjectResponse, len(objects))
	for i, obj := range objects {
		responses[i] = dto.NewObjectResponse(obj)
	}

	return responses, nil
}

func (s *ObjectService) DeleteObject(ctx context.Context, functionID, name string) error {
	return s.objectRepo.Delete(ctx, functionID, name)
}
