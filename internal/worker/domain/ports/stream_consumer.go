package ports

import (
	"context"
	"faas/internal/features/executions/domain/entity"
)

type StreamConsumer interface {
	Subscribe(handler func(ctx context.Context, execution *entity.Execution) error) Worker
}

type Worker interface {
	Stop() error
}
