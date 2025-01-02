package ports

import (
	"context"
	"faas/internal/features/functions/domain/entity"
	secretEntity "faas/internal/features/secrets/domain/entity"
)

type FunctionRepository interface {
	GetByID(ctx context.Context, id string) (*entity.Function, error)
}

type SecretRepository interface {
	GetByName(ctx context.Context, userID, name string) (*secretEntity.Secret, error)
}
