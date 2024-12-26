package repository

import (
	"context"
	"encoding/json"

	"faas/internal/features/functions/domain/entity"
	"faas/internal/shared/infrastructure/nats"
)

type NatsFunctionRepository struct {
	kv nats.KeyValue
}

func NewNatsFunctionRepository(js nats.JetStreamContext) (*NatsFunctionRepository, error) {
	natsKV, err := js.KeyValue(nats.FUNCTIONS_BUCKET)
	if err != nil {
		return nil, err
	}

	return &NatsFunctionRepository{
		kv: nats.NewKeyValueAdapter(natsKV),
	}, nil
}

func (r *NatsFunctionRepository) Save(ctx context.Context, function *entity.Function) error {
	data, err := json.Marshal(function)
	if err != nil {
		return err
	}

	_, err = r.kv.Put(function.ID, data)
	return err
}

func (r *NatsFunctionRepository) ListByUserID(ctx context.Context, userID string) ([]*entity.Function, error) {
	keys, err := r.kv.Keys()
	if err != nil {
		return nil, err
	}

	var functions []*entity.Function
	for _, key := range keys {
		entry, err := r.kv.Get(key)
		if err != nil {
			continue
		}

		var function entity.Function
		if err := json.Unmarshal(entry.Value(), &function); err != nil {
			continue
		}

		if function.UserID == userID {
			functions = append(functions, &function)
		}
	}

	return functions, nil
}

func (r *NatsFunctionRepository) Delete(ctx context.Context, id string) error {
	return r.kv.Delete(id)
}

func (r *NatsFunctionRepository) GetByID(ctx context.Context, id string) (*entity.Function, error) {
	entry, err := r.kv.Get(id)
	if err != nil {
		return nil, err
	}

	var function entity.Function
	if err := json.Unmarshal(entry.Value(), &function); err != nil {
		return nil, err
	}

	return &function, nil
}
