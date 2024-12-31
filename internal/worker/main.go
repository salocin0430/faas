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

	// Conectar a NATS
	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	// Crear JetStream Context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("Failed to create JetStream context:", err)
	}

	// Inicializar componentes
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

	// Crear servicio
	executionService := service.NewExecutionService(
		containerManager,
		executionRepo,
	)

	log.Println("Starting worker...")

	// Configurar consumer
	worker := streamConsumer.Subscribe(executionService.ProcessExecution)
	log.Println("Subscribed to executions.pending")

	// Manejar shutdown gracefully
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Canal para señales de sistema
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Esperar señal o cancelación de contexto
	select {
	case <-sigChan:
		log.Println("Shutting down worker...")
	case <-ctx.Done():
		log.Println("Context cancelled, shutting down...")
	}

	// Cleanup
	worker.Stop()
}
