package ports

import (
	"context"
	"faas/internal/features/secrets/domain/entity"
)

type SecretRepository interface {
	Create(ctx context.Context, secret *entity.Secret) error
	GetByID(ctx context.Context, userID, secretID string) (*entity.Secret, error)
	GetByName(ctx context.Context, userID, name string) (*entity.Secret, error)
	List(ctx context.Context, userID string) ([]*entity.Secret, error)
	Update(ctx context.Context, secret *entity.Secret) error
	Delete(ctx context.Context, userID, secretID string) error
}
