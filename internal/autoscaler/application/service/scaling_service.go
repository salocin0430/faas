package service

import (
	"context"
	"faas/internal/autoscaler/domain/entity"
	"faas/internal/autoscaler/domain/repository"
	"time"
)

type ScalingService struct {
	metricsRepo repository.MetricsRepository
	policy      *entity.ScalingPolicy
	lastScaled  time.Time
}

func NewScalingService(repo repository.MetricsRepository, policy *entity.ScalingPolicy) *ScalingService {
	return &ScalingService{
		metricsRepo: repo,
		policy:      policy,
	}
}

func (s *ScalingService) EvaluateScaling(ctx context.Context) (int, error) {
	// Si estamos en periodo de cooldown, no escalar
	if time.Since(s.lastScaled) < s.policy.CooldownPeriod {
		return 0, nil
	}

	metrics, err := s.metricsRepo.GetSystemMetrics(ctx)
	if err != nil {
		return 0, err
	}

	// Calcular promedio de métricas
	var totalCPU, totalMemory float64
	for _, metric := range metrics {
		totalCPU += metric.CPU
		totalMemory += metric.Memory
	}
	avgCPU := totalCPU / float64(len(metrics))
	avgMemory := totalMemory / float64(len(metrics))

	// Decidir si escalar
	if avgCPU > s.policy.CPUThreshold || avgMemory > s.policy.MemoryThreshold {
		s.lastScaled = time.Now()
		return s.policy.ScaleUpStep, nil
	}

	return 0, nil
}
