package repository

import (
	"context"
	"encoding/json"
	"faas/internal/features/users/domain/entity"
	"faas/internal/shared/infrastructure/nats"
	"fmt"
	"log"
)

type NatsUserRepository struct {
	kv nats.KeyValue
}

func NewNatsUserRepository(js nats.JetStreamContext) (*NatsUserRepository, error) {
	natsKV, err := js.KeyValue(nats.USERS_BUCKET)
	if err != nil {
		return nil, err
	}

	return &NatsUserRepository{
		kv: nats.NewKeyValueAdapter(natsKV),
	}, nil
}

func (r *NatsUserRepository) Create(ctx context.Context, user *entity.User) error {
	log.Printf("Creating user with hash: %s", user.Password)

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	log.Printf("Serialized user data: %s", string(data))

	_, err = r.kv.Put(user.ID, data)
	return err
}

func (r *NatsUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	entry, err := r.kv.Get(id)
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := json.Unmarshal(entry.Value(), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *NatsUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	users, err := r.List(ctx)
	if err != nil {
		log.Printf("Error listing users: %v", err)
		return nil, err
	}

	for _, user := range users {
		log.Printf("Found user with hash: %s", user.Password)
		if user.Username == username {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

func (r *NatsUserRepository) Delete(ctx context.Context, id string) error {
	return r.kv.Delete(id)
}

func (r *NatsUserRepository) List(ctx context.Context) ([]*entity.User, error) {
	keys, err := r.kv.Keys()
	if err != nil {
		return nil, err
	}

	var users []*entity.User
	for _, key := range keys {
		entry, err := r.kv.Get(key)
		if err != nil {
			continue
		}

		var user entity.User
		if err := json.Unmarshal(entry.Value(), &user); err != nil {
			continue
		}

		users = append(users, &user)
	}

	return users, nil
}
