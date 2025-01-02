package repository

import (
	"context"
	"encoding/json"
	"faas/internal/features/secrets/domain/entity"
	"faas/internal/shared/domain/errors"
	"faas/internal/shared/infrastructure/nats"
	"fmt"
	"strings"
)

type KVSecretRepository struct {
	kv nats.KeyValue
}

func NewKVSecretRepository(js nats.JetStreamContext) (*KVSecretRepository, error) {
	kv, err := js.KeyValue(nats.SECRETS_BUCKET)
	if err != nil {
		// If bucket already exists, try to get it
		kv, err = js.KeyValue(nats.SECRETS_BUCKET)
		if err != nil {
			return nil, err
		}
	}

	return &KVSecretRepository{
		kv: nats.NewKeyValueAdapter(kv),
	}, nil
}

func (r *KVSecretRepository) Create(ctx context.Context, secret *entity.Secret) error {
	data, err := json.Marshal(secret)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s", secret.UserID, secret.ID)
	_, err = r.kv.Put(key, data)
	return err
}

func (r *KVSecretRepository) GetByID(ctx context.Context, userID, secretID string) (*entity.Secret, error) {
	key := fmt.Sprintf("%s/%s", userID, secretID)
	entry, err := r.kv.Get(key)
	if err != nil {
		if err.Error() == "nats: key not found" {
			return nil, errors.NewAppError("secret_not_found", "Secret not found")
		}
		return nil, err
	}

	var secret entity.Secret
	if err := json.Unmarshal(entry.Value(), &secret); err != nil {
		return nil, err
	}

	return &secret, nil
}

func (r *KVSecretRepository) GetByName(ctx context.Context, userID, name string) (*entity.Secret, error) {
	// List all secrets for user and find by name
	secrets, err := r.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, secret := range secrets {
		if secret.Name == name {
			return secret, nil
		}
	}

	return nil, errors.NewAppError("secret_not_found", "Secret not found")
}

func (r *KVSecretRepository) List(ctx context.Context, userID string) ([]*entity.Secret, error) {
	prefix := fmt.Sprintf("%s/", userID)
	keys, err := r.kv.Keys()
	if err != nil {
		return nil, err
	}

	secrets := make([]*entity.Secret, 0, len(keys))
	for _, key := range keys {
		if strings.HasPrefix(key, prefix) {
			entry, err := r.kv.Get(key)
			if err != nil {
				continue
			}

			var secret entity.Secret
			if err := json.Unmarshal(entry.Value(), &secret); err != nil {
				continue
			}

			secrets = append(secrets, &secret)
		}
	}

	return secrets, nil
}

func (r *KVSecretRepository) Update(ctx context.Context, secret *entity.Secret) error {
	return r.Create(ctx, secret) // Same as create since we're using Put
}

func (r *KVSecretRepository) Delete(ctx context.Context, userID, secretID string) error {
	key := fmt.Sprintf("%s/%s", userID, secretID)
	return r.kv.Delete(key)
}
