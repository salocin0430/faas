package nats

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"faas/internal/features/scaling/domain/entity"
	"faas/internal/features/scaling/domain/service"
)

type NatsAutoScaler struct {
	natsClient       *nats.Conn
	metricsCollector entity.MetricsCollector
	policy           service.ScalingPolicy
	config           service.AutoScalerConfig
	currentWorkers   int
}

func NewNatsAutoScaler(nc *nats.Conn, config service.AutoScalerConfig) service.AutoScaler {
	return &NatsAutoScaler{
		natsClient:       nc,
		metricsCollector: NewNatsMetricsCollector(nc),
		policy:           &service.DefaultScalingPolicy{},
		config:           config,
		currentWorkers:   config.MinWorkers,
	}
}

func (as *NatsAutoScaler) Start(ctx context.Context) error {
	log.Printf("Starting AutoScaler (min=%d, max=%d)", as.config.MinWorkers, as.config.MaxWorkers)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			metrics := as.metricsCollector.Collect()
			as.evaluateAndScale(metrics)
		}
	}
}

func (as *NatsAutoScaler) Stop() error {
	if as.natsClient != nil {
		as.natsClient.Close()
	}
	return nil
}

func (as *NatsAutoScaler) GetCurrentWorkers() int {
	return as.currentWorkers
}

func (as *NatsAutoScaler) evaluateAndScale(metrics entity.Metrics) {
	decision := as.policy.EvaluateScaling(
		metrics,
		as.currentWorkers,
		as.config.MinWorkers,
		as.config.MaxWorkers,
	)

	if decision.Action != "none" {
		as.currentWorkers = decision.Target
		as.publishScalingEvent(decision)
	}
}

func (as *NatsAutoScaler) publishScalingEvent(decision service.ScalingDecision) {
	event := map[string]interface{}{
		"action":  decision.Action,
		"workers": decision.Target,
		"time":    time.Now().Unix(),
	}

	data, _ := json.Marshal(event)
	as.natsClient.Publish("scaling_events", data)
	log.Printf("Published scaling event: %s to %d workers", decision.Action, decision.Target)
}
