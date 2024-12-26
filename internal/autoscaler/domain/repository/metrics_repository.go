package repository

import (
	"context"
	"faas/internal/autoscaler/domain/entity"
)

type MetricsRepository interface {
	SaveMetric(ctx context.Context, metric *entity.Metric) error
	GetWorkerMetrics(ctx context.Context, workerID string) ([]*entity.Metric, error)
	GetSystemMetrics(ctx context.Context) ([]*entity.Metric, error)
}
