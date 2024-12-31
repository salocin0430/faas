package nats

import (
	"context"
	"encoding/json"
	"faas/internal/features/executions/domain/entity"
	"faas/internal/shared/infrastructure/nats"
	"faas/internal/worker/domain/ports"
)

type NatsExecutionRepository struct {
	kv nats.KeyValue
}

func NewExecutionRepository(js nats.JetStreamContext) (ports.ExecutionRepository, error) {
	kv, err := js.KeyValue(nats.EXECUTIONS_BUCKET)
	if err != nil {
		return nil, err
	}
	return &NatsExecutionRepository{
		kv: nats.NewKeyValueAdapter(kv),
	}, nil
}

func (r *NatsExecutionRepository) UpdateExecution(ctx context.Context, execution *entity.Execution) error {
	data, err := json.Marshal(execution)
	if err != nil {
		return err
	}
	_, err = r.kv.Put(execution.ID, data)
	return err
}
