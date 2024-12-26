package repository

import (
	"context"
	"encoding/json"
	"faas/internal/features/executions/domain/entity"
	"faas/internal/shared/infrastructure/nats"
)

type NatsExecutionRepository struct {
	kv nats.KeyValue
}

func NewNatsExecutionRepository(js nats.JetStreamContext) (*NatsExecutionRepository, error) {
	kv, err := js.KeyValue(nats.EXECUTIONS_BUCKET)
	if err != nil {
		return nil, err
	}
	return &NatsExecutionRepository{kv: nats.NewKeyValueAdapter(kv)}, nil
}

func (r *NatsExecutionRepository) Save(ctx context.Context, execution *entity.Execution) error {
	data, err := json.Marshal(execution)
	if err != nil {
		return err
	}
	_, err = r.kv.Put(execution.ID, data)
	return err
}

func (r *NatsExecutionRepository) GetByID(ctx context.Context, id string) (*entity.Execution, error) {
	data, err := r.kv.Get(id)
	if err != nil {
		return nil, err
	}

	var execution entity.Execution
	if err := json.Unmarshal(data.Value(), &execution); err != nil {
		return nil, err
	}
	return &execution, nil
}

func (r *NatsExecutionRepository) ListByUserID(ctx context.Context, userID string) ([]*entity.Execution, error) {
	keys, err := r.kv.Keys()
	if err != nil {
		return nil, err
	}

	var executions []*entity.Execution
	for _, key := range keys {
		entry, err := r.kv.Get(key)
		if err != nil {
			continue
		}

		var execution entity.Execution
		if err := json.Unmarshal(entry.Value(), &execution); err != nil {
			continue
		}

		if execution.UserID == userID {
			executions = append(executions, &execution)
		}
	}

	return executions, nil
}

func (r *NatsExecutionRepository) Update(ctx context.Context, execution *entity.Execution) error {
	return r.Save(ctx, execution)
}
