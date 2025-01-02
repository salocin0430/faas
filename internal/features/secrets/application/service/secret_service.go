package service

import (
	"context"
	"faas/internal/features/secrets/application/dto"
	"faas/internal/features/secrets/domain/entity"
	"faas/internal/features/secrets/domain/ports"
	"faas/internal/shared/domain/errors"
	"time"

	"github.com/google/uuid"
)

type SecretService struct {
	secretRepo ports.SecretRepository
}

func NewSecretService(secretRepo ports.SecretRepository) *SecretService {
	return &SecretService{
		secretRepo: secretRepo,
	}
}

func (s *SecretService) CreateSecret(ctx context.Context, userID string, req *dto.CreateSecretRequest) (*dto.SecretResponse, error) {
	// Check if secret with same name exists
	if _, err := s.secretRepo.GetByName(ctx, userID, req.Name); err == nil {
		return nil, errors.NewAppError("secret_with_this_name_already_exists", "Secret with this name already exists")
	}

	secret := &entity.Secret{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      req.Name,
		Value:     req.Value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.secretRepo.Create(ctx, secret); err != nil {
		return nil, err
	}

	return dto.ToSecretResponse(secret), nil
}

func (s *SecretService) GetSecret(ctx context.Context, userID, secretID string) (*dto.SecretResponse, error) {
	secret, err := s.secretRepo.GetByID(ctx, userID, secretID)
	if err != nil {
		return nil, err
	}

	return dto.ToSecretResponse(secret), nil
}

func (s *SecretService) ListSecrets(ctx context.Context, userID string) ([]*dto.SecretResponse, error) {
	secrets, err := s.secretRepo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.ToSecretResponseList(secrets), nil
}

func (s *SecretService) UpdateSecret(ctx context.Context, userID, secretID string, req *dto.UpdateSecretRequest) (*dto.SecretResponse, error) {
	secret, err := s.secretRepo.GetByID(ctx, userID, secretID)
	if err != nil {
		return nil, err
	}

	secret.Value = req.Value
	secret.UpdatedAt = time.Now()

	if err := s.secretRepo.Update(ctx, secret); err != nil {
		return nil, err
	}

	return dto.ToSecretResponse(secret), nil
}

func (s *SecretService) DeleteSecret(ctx context.Context, userID, secretID string) error {
	return s.secretRepo.Delete(ctx, userID, secretID)
}
