package dto

import (
	"faas/internal/features/secrets/domain/entity"
	"time"
)

type SecretResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToSecretResponse(secret *entity.Secret) *SecretResponse {
	return &SecretResponse{
		ID:        secret.ID,
		Name:      secret.Name,
		CreatedAt: secret.CreatedAt,
		UpdatedAt: secret.UpdatedAt,
	}
}

func ToSecretResponseList(secrets []*entity.Secret) []*SecretResponse {
	responses := make([]*SecretResponse, len(secrets))
	for i, secret := range secrets {
		responses[i] = ToSecretResponse(secret)
	}
	return responses
}
