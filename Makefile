.PHONY: build run test down clean

# Image names
API_IMAGE=faas-api
WORKER_IMAGE=faas-worker
AUTOSCALER_IMAGE=faas-autoscaler

# Build Docker images
build:
	docker build -t $(API_IMAGE) -f deployments/docker/Dockerfile.api .
	docker build -t $(WORKER_IMAGE) -f deployments/docker/Dockerfile.worker .
# docker build -t $(AUTOSCALER_IMAGE) -f deployments/docker/Dockerfile.autoscaler .

# Run with docker-compose
run:
	docker-compose up

# Run tests
test:
	docker-compose run --rm api go test ./...

# Stop services
down:
	docker-compose down

# Clean volumes and cache
clean:
	docker-compose down -v
	docker system prune -f

# Help
help:
	@echo "Available commands:"
	@echo "  build  - Build Docker images"
	@echo "  run    - Run services with docker-compose"
	@echo "  test   - Run tests"
	@echo "  down   - Stop services"
	@echo "  clean  - Clean volumes and cache"