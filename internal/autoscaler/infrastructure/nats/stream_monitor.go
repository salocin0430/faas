package nats

import (
	"faas/internal/autoscaler/domain/ports"

	"github.com/nats-io/nats.go"
)

type NatsStreamMonitor struct {
	js nats.JetStreamContext
}

func NewStreamMonitor(js nats.JetStreamContext) ports.StreamMonitor {
	return &NatsStreamMonitor{js: js}
}

func (m *NatsStreamMonitor) GetPendingMessages() (int, error) {
	// Obtener info del stream
	stream, err := m.js.StreamInfo("EXECUTIONS")
	if err != nil {
		return 0, err
	}

	// Retornar mensajes pendientes
	return int(stream.State.Msgs), nil
}