package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"faas/internal/autoscaler/application/service"
	"faas/internal/autoscaler/domain/config"
	"faas/internal/autoscaler/infrastructure/docker"
	"faas/internal/autoscaler/infrastructure/nats"

	natspkg "github.com/nats-io/nats.go"
)

func validateConfig(cfg *config.AutoscalerConfig) error {
	if cfg.MinWorkers < 2 {
		return fmt.Errorf("MinWorkers must be at least 2")
	}
	if cfg.MaxWorkers < cfg.MinWorkers {
		return fmt.Errorf("MaxWorkers must be greater than MinWorkers")
	}
	if _, err := time.ParseDuration(cfg.CheckInterval); err != nil {
		return fmt.Errorf("invalid CheckInterval: %v", err)
	}
	if _, err := time.ParseDuration(cfg.CooldownPeriod); err != nil {
		return fmt.Errorf("invalid CooldownPeriod: %v", err)
	}
	return nil
}

func main() {
	// Load configuration
	cfg := &config.AutoscalerConfig{
		MinWorkers:         2,
		MaxWorkers:         10,
		ScaleUpThreshold:   3,
		ScaleDownThreshold: 2,
		CheckInterval:      "30s",
		CooldownPeriod:     "30s",
	}

	if err := validateConfig(cfg); err != nil {
		log.Fatal("Invalid configuration:", err)
	}

	// Connect to NATS
	nc, err := natspkg.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("Failed to create JetStream context:", err)
	}

	// Initialize components
	streamMonitor := nats.NewStreamMonitor(js)
	scaler, err := docker.NewDockerScaler("worker")
	if err != nil {
		log.Fatal("Failed to create docker scaler:", err)
	}

	// Create service
	autoscaler := service.NewAutoscalerService(scaler, streamMonitor, cfg)

	// Handle termination signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Shutting down autoscaler...")
		cancel()
	}()

	// Start autoscaler
	log.Println("Starting autoscaler...")
	if err := autoscaler.Start(ctx); err != nil {
		log.Fatal("Autoscaler failed:", err)
	}
}
