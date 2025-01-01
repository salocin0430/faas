package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"faas/internal/autoscaler/domain/config"
	"faas/internal/autoscaler/domain/ports"
)

type AutoscalerService struct {
	scaler        ports.Scaler
	streamMonitor ports.StreamMonitor
	config        *config.AutoscalerConfig
	lastScaleTime time.Time
	cooldown      time.Duration
}

func NewAutoscalerService(scaler ports.Scaler, monitor ports.StreamMonitor, config *config.AutoscalerConfig) *AutoscalerService {
	cooldown, _ := time.ParseDuration(config.CooldownPeriod)
	return &AutoscalerService{
		scaler:        scaler,
		streamMonitor: monitor,
		config:        config,
		lastScaleTime: time.Now(),
		cooldown:      cooldown,
	}
}

func (s *AutoscalerService) Start(ctx context.Context) error {
	checkInterval, _ := time.ParseDuration(s.config.CheckInterval)
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := s.check(); err != nil {
				log.Printf("Error during scaling check: %v", err)
			}
		}
	}
}

func (s *AutoscalerService) check() error {
	// Obtener métricas
	pendingMsgs, err := s.streamMonitor.GetPendingMessages()
	if err != nil {
		return err
	}

	workers, err := s.scaler.GetCurrentWorkers()
	if err != nil {
		return err
	}

	// Calcular métricas
	msgsPerWorker := float64(pendingMsgs) / float64(workers)

	log.Printf("Status Check - Pending Messages: %d, Active Workers: %d, Messages per Worker: %.2f",
		pendingMsgs, workers, msgsPerWorker)

	// Decidir escalado
	if msgsPerWorker > float64(s.config.ScaleUpThreshold) {
		log.Printf("Scaling UP - Too many messages per worker (%.2f > %d)",
			msgsPerWorker, s.config.ScaleUpThreshold)
		if err := s.scaleUp(); err != nil {
			return fmt.Errorf("scale up failed: %v", err)
		}
	} else if msgsPerWorker < float64(s.config.ScaleDownThreshold) {
		log.Printf("Scaling DOWN - Few messages per worker (%.2f < %d)",
			msgsPerWorker, s.config.ScaleDownThreshold)
		if err := s.scaleDown(); err != nil {
			return fmt.Errorf("scale down failed: %v", err)
		}
	} else {
		log.Printf("No scaling needed - Current load is optimal")
	}

	return nil
}

func (s *AutoscalerService) scaleUp() error {
	if time.Since(s.lastScaleTime) < s.cooldown {
		return nil
	}
	workers, err := s.scaler.GetCurrentWorkers()
	if err != nil {
		return err
	}
	if workers >= s.config.MaxWorkers {
		return nil
	}
	s.lastScaleTime = time.Now()
	return s.scaler.ScaleUp(1)
}

func (s *AutoscalerService) scaleDown() error {
	if time.Since(s.lastScaleTime) < s.cooldown {
		return nil
	}
	workers, err := s.scaler.GetCurrentWorkers()
	if err != nil {
		return err
	}
	if workers <= s.config.MinWorkers {
		return nil
	}
	s.lastScaleTime = time.Now()
	return s.scaler.ScaleDown(1)
}
