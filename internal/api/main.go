package main

import (
	"log"

	"faas/internal/shared/infrastructure/config"
	"faas/internal/shared/infrastructure/nats"

	"github.com/gin-gonic/gin"

	funcService "faas/internal/features/functions/application/service"
	funcRepo "faas/internal/features/functions/infrastructure/repository"
	funcHttp "faas/internal/features/functions/interfaces/http"

	objService "faas/internal/features/function_objects/application/service"
	objRepo "faas/internal/features/function_objects/infrastructure/nats"
	objHttp "faas/internal/features/function_objects/interfaces/http"

	userService "faas/internal/features/users/application/service"
	userRepo "faas/internal/features/users/infrastructure/repository"
	userHttp "faas/internal/features/users/interfaces/http"

	execService "faas/internal/features/executions/application/service"
	execRepo "faas/internal/features/executions/infrastructure/repository"
	execHttp "faas/internal/features/executions/interfaces/http"

	secretService "faas/internal/features/secrets/application/service"
	secretRepo "faas/internal/features/secrets/infrastructure/repository"
	secretHttp "faas/internal/features/secrets/interfaces/http"
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

	// Create buckets in NATS
	if err := nats.CreateBuckets(js); err != nil {
		log.Fatal("Failed to create NATS buckets:", err)
	}

	// Create streams in NATS
	if err := nats.CreateStreams(js); err != nil {
		log.Fatal("Failed to create NATS streams:", err)
	}

	// Initialize repositories
	functionRepo, err := funcRepo.NewNatsFunctionRepository(js)
	if err != nil {
		log.Fatal(err)
	}

	userRepo, err := userRepo.NewNatsUserRepository(js)
	if err != nil {
		log.Fatal(err)
	}

	executionRepo, err := execRepo.NewNatsExecutionRepository(js)
	if err != nil {
		log.Fatal(err)
	}

	objectRepo, err := objRepo.NewNatsObjectRepository(js)
	if err != nil {
		log.Fatal(err)
	}

	secretRepo, err := secretRepo.NewKVSecretRepository(js)
	if err != nil {
		log.Fatal(err)
	}

	// Stream repository
	execStreamRepo := execRepo.NewNatsExecutionStreamRepository(js)

	// Initialize services
	funcService := funcService.NewFunctionService(functionRepo)
	userService := userService.NewUserService(userRepo, cfg)
	executionService := execService.NewExecutionService(executionRepo, execStreamRepo, functionRepo, cfg)
	objectService := objService.NewObjectService(objectRepo)
	secretService := secretService.NewSecretService(secretRepo)
	// Initialize handlers
	functionHandler := funcHttp.NewFunctionHandler(funcService)
	userHandler := userHttp.NewUserHandler(userService)
	executionHandler := execHttp.NewExecutionHandler(executionService)
	objectHandler := objHttp.NewObjectHandler(objectService)
	secretHandler := secretHttp.NewSecretHandler(secretService)
	// Initialize Gin
	r := gin.Default()

	// Setup routes
	funcHttp.SetupFunctionRoutes(r, functionHandler, cfg.JWTSecret)
	userHttp.SetupUserRoutes(r, userHandler)
	execHttp.SetupExecutionRoutes(r, executionHandler, cfg.JWTSecret)
	objHttp.SetupObjectRoutes(r, objectHandler)
	secretHttp.SetupSecretRoutes(r, secretHandler, cfg.JWTSecret)
	// Start server
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
