package nats

import (
	"encoding/json"
	"log"
	"time"

	"faas/internal/features/scaling/domain/entity"

	"github.com/nats-io/nats.go"
)

type NatsMetricsCollector struct {
	natsClient *nats.Conn
}

func NewNatsMetricsCollector(nc *nats.Conn) entity.MetricsCollector {
	return &NatsMetricsCollector{
		natsClient: nc,
	}
}

func (mc *NatsMetricsCollector) Collect() entity.Metrics {
	// Obtener estadísticas de NATS
	//stats := mc.natsClient.Stats()

	metrics := entity.Metrics{
		//QueueLength:   len(stats.PendingMsgs),
		//ActiveWorkers: len(stats.SubChannels),
		Timestamp: time.Now().Unix(),
	}

	// Calcular tiempo promedio de procesamiento
	metrics.ProcessingTime = mc.calculateAverageProcessingTime()

	log.Printf("Collected metrics: %+v", metrics)
	return metrics
}

func (mc *NatsMetricsCollector) PublishMetrics(metrics entity.Metrics) error {
	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	return mc.natsClient.Publish("metrics", data)
}

func (mc *NatsMetricsCollector) calculateAverageProcessingTime() float64 {
	// Implementación simple: podríamos mejorarla guardando históricos
	// Por ahora retorna un valor fijo para pruebas
	return 1.0
}
