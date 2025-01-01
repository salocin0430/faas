package nats

import (
	"context"
	"encoding/json"
	"faas/internal/autoscaler/domain/entity"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsMetricsRepository struct {
	kv nats.KeyValue
}

func NewNatsMetricsRepository(js nats.JetStreamContext) (*NatsMetricsRepository, error) {
	kv, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:  "metrics",
		History: 5,
	})
	if err != nil {
		return nil, err
	}
	return &NatsMetricsRepository{kv: kv}, nil
}

func (r *NatsMetricsRepository) SaveMetric(ctx context.Context, metric *entity.Metric) error {
	data, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", metric.WorkerID, metric.CollectedAt.Format(time.RFC3339))
	_, err = r.kv.Put(key, data)
	return err
}

func (r *NatsMetricsRepository) GetWorkerMetrics(ctx context.Context, workerID string) ([]*entity.Metric, error) {
	// Implement search for metrics by worker
	return nil, nil
}

func (r *NatsMetricsRepository) GetSystemMetrics(ctx context.Context) ([]*entity.Metric, error) {
	// Implement search for all system metrics
	return nil, nil
}
