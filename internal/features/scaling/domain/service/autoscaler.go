package service

import (
	"context"

	"faas/internal/features/scaling/domain/entity"
)

const (
	HIGH_THRESHOLD = 100 // Número de mensajes en cola para escalar arriba
	LOW_THRESHOLD  = 10  // Número de mensajes en cola para escalar abajo
)

// AutoScaler define la interfaz del servicio de auto-escalado
type AutoScaler interface {
	Start(ctx context.Context) error
	Stop() error
	GetCurrentWorkers() int
}

// AutoScalerConfig contiene la configuración del auto-escalador
type AutoScalerConfig struct {
	MinWorkers int
	MaxWorkers int
}

// ScalingPolicy define la política de escalado
type ScalingPolicy interface {
	EvaluateScaling(metrics entity.Metrics, currentWorkers, minWorkers, maxWorkers int) ScalingDecision
}

// ScalingDecision representa una decisión de escalado
type ScalingDecision struct {
	Action string // "scale_up", "scale_down", "none"
	Target int    // Número objetivo de workers
}

// DefaultScalingPolicy implementa la política de escalado por defecto
type DefaultScalingPolicy struct{}

func (p *DefaultScalingPolicy) EvaluateScaling(metrics entity.Metrics, currentWorkers, minWorkers, maxWorkers int) ScalingDecision {
	if metrics.QueueLength > HIGH_THRESHOLD && currentWorkers < maxWorkers {
		return ScalingDecision{
			Action: "scale_up",
			Target: currentWorkers + 1,
		}
	}

	if metrics.QueueLength < LOW_THRESHOLD && currentWorkers > minWorkers {
		return ScalingDecision{
			Action: "scale_down",
			Target: currentWorkers - 1,
		}
	}

	return ScalingDecision{
		Action: "none",
		Target: currentWorkers,
	}
}
