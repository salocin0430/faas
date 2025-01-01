package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	funcRepo "faas/internal/features/functions/infrastructure/repository"
	"faas/internal/shared/infrastructure/config"
	"faas/internal/shared/infrastructure/nats"
	"faas/internal/worker/application/service"
	"faas/internal/worker/infrastructure/docker"
	workerNats "faas/internal/worker/infrastructure/nats"
)

func main() {
	cfg := config.LoadConfig()

	// Connect to NATS
	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	// Create JetStream Context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("Failed to create JetStream context:", err)
	}

	// Initialize components
	functionRepo, err := funcRepo.NewNatsFunctionRepository(js)
	if err != nil {
		log.Fatal("Failed to create function repository:", err)
	}

	containerManager, err := docker.NewContainerManager(functionRepo)
	if err != nil {
		log.Fatal("Failed to create container manager:", err)
	}

	executionRepo, err := workerNats.NewExecutionRepository(js)
	if err != nil {
		log.Fatal("Failed to create execution repository:", err)
	}

	streamConsumer := workerNats.NewStreamConsumer(js)

	// Create service
	executionService := service.NewExecutionService(
		containerManager,
		executionRepo,
	)

	log.Println("Starting worker...")

	// Configure consumer
	worker := streamConsumer.Subscribe(executionService.ProcessExecution)
	log.Println("Subscribed to executions.pending")

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel for system signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal or context cancellation
	select {
	case <-sigChan:
		log.Println("Shutting down worker...")
	case <-ctx.Done():
		log.Println("Context cancelled, shutting down...")
	}

	// Cleanup
	worker.Stop()
}
