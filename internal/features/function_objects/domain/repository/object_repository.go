package repository

import (
	"context"
	"faas/internal/features/function_objects/domain/entity"
)

type ObjectRepository interface {
	Save(ctx context.Context, obj *entity.FunctionObject, data []byte) error
	Get(ctx context.Context, functionID, objectName string) (*entity.FunctionObject, []byte, error)
	List(ctx context.Context, functionID string) ([]*entity.FunctionObject, error)
	Delete(ctx context.Context, functionID, objectName string) error
}
