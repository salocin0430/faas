package repository

import (
	"encoding/json"
	"faas/internal/features/executions/domain/entity"
	"faas/internal/shared/infrastructure/nats"
)

type NatsExecutionStreamRepository struct {
	js nats.JetStreamContext
}

func NewNatsExecutionStreamRepository(js nats.JetStreamContext) *NatsExecutionStreamRepository {
	return &NatsExecutionStreamRepository{js: js}
}

func (r *NatsExecutionStreamRepository) PublishPending(execution *entity.Execution) error {
	data, err := json.Marshal(execution)
	if err != nil {
		return err
	}
	_, err = r.js.Publish("executions.pending", data)
	return err
}
