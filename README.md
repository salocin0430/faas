# FaaS Platform

A Function as a Service (FaaS) platform that allows users to deploy and execute containerized functions with support for secrets management and file processing.

## Architecture

### Overview
The platform follows a Clean Architecture pattern with the following layers:
- Domain Layer: Core business logic and entities
- Application Layer: Use cases and service orchestration
- Infrastructure Layer: External implementations (NATS, Docker, etc.)
- Interfaces Layer: HTTP handlers and routes

### Key Components
1. **API Service**
   - Handles HTTP requests
   - User authentication and authorization
   - Function and execution management
   - Object storage operations
   - Secrets management

2. **Worker Service**
   - Processes function executions
   - Container management with Docker
   - Handles function input/output
   - Manages execution timeouts

3. **Storage Layer**
   - NATS JetStream for event streaming
   - NATS KeyValue for data persistence
   - Object storage for function files

4. **Gateway**
   - Apache APISIX for API gateway
   - Route management
   - Authentication middleware
   - CORS and security policies

### Technologies Used
- Go 1.21+
- NATS 2.0+
- Docker API
- Apache APISIX
- JWT for authentication

## Features

### Core Features
1. **Function Management**
   - Create, list, and delete functions
   - Docker image-based functions
   - Function metadata management

2. **Execution Engine**
   - Asynchronous function execution
   - Input/output handling
   - Execution status tracking
   - Timeout management

3. **Object Storage**
   - File upload and download
   - Object lifecycle management
   - Function-specific object storage

4. **Secrets Management**
   - Secure secret storage
   - Per-user secret scoping
   - Runtime secret injection

5. **User Management**
   - User registration and authentication
   - Role-based access control
   - JWT-based session management

### Security Features
- JWT-based authentication
- Per-user resource isolation
- Secure secret management
- Network isolation for functions

## Configuration

### Environment Variables
```env
# Server Configuration
SERVER_ADDRESS=":8080"
JWT_SECRET="your-super-secret-key-for-development"
CONSUMER_KEY="faasapp-key"
MAX_CONCURRENT_EXECUTIONS="10"

# NATS Configuration
NATS_URL="nats://localhost:4222"

# Docker Configuration
NETWORK_NAME="apisix"
API_BASE_URL="http://api:8080/api/function-objects"
```

### Docker Compose Setup
```yaml
version: '3.8'

services:
  apisix:
    image: apache/apisix:3.11.0-debian
    restart: always
    volumes:
      - ./deployments/apisix/config.yaml:/usr/local/apisix/conf/config.yaml:ro  
      - ./deployments/apisix/apisix.yaml:/usr/local/apisix/conf/apisix.yaml:ro
    ports:
      - "9080:9080"
    environment:
      - GATEWAY_PORT=9080
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9080/apisix/admin/health"]
      interval: 5s
      timeout: 3s
      retries: 5
    networks:
      - apisix

  backend:
    image: nginx:1.25-alpine
    restart: always
    volumes:
      - ./deployments/upstream/backend.conf:/etc/nginx/nginx.conf:ro
    ports:
      - "9081:80/tcp"
    environment:
      - NGINX_PORT=80    
    networks:
      - apisix

  # Services
  nats:
    image: nats:latest
    command: ["--js"]
    ports:
      - "4222:4222"
      - "8222:8222"
    volumes:
      - nats-data:/data/nats-server/jetstream
    networks:
      - apisix

  api:
    build:
      context: .
      dockerfile: deployments/docker/Dockerfile.api
    expose:
      - "8080"
    environment:
      - NATS_URL=nats://nats:4222
      - SERVER_ADDRESS=:8080
      - JWT_SECRET=your-super-secret-key-for-development  
    depends_on:
      - nats
      - apisix
    networks:
      - apisix

  worker:
    build:
      context: .
      dockerfile: deployments/docker/Dockerfile.worker
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    group_add:
      - ${DOCKER_GID:-999}
    environment:
      - NATS_URL=nats://nats:4222
    depends_on:
      - nats
      - api
    networks:
      - apisix
    deploy:
      replicas: 2

  autoscaler:
    build:
      context: .
      dockerfile: deployments/docker/Dockerfile.autoscaler
    environment:
      - NATS_URL=nats://nats:4222
      - MIN_WORKERS=2
      - MAX_WORKERS=10
      - SCALE_UP_THRESHOLD=3
      - SCALE_DOWN_THRESHOLD=2
      - CHECK_INTERVAL=30s
      - COOLDOWN_PERIOD=30s
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock  
    depends_on:
      - nats
      - worker
    networks:
      - apisix

volumes:
  nats-data:
  apisix_data:

networks:
  apisix:  
    name: apisix
    driver: bridge
```

## Getting Started

1. Clone the repository
```bash
git clone https://github.com/salocin0430/faas.git
```

2. Set up environment variables
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Start the services
## Build and Run

The project uses Make for common operations:

```bash
# Build Docker images
make build

# Run all services
make run

# Stop services
make down

# Run tests
make test

# Clean up everything
make clean
```

For more details:
```bash
make help
```

4. Create a user and get a token
```bash
curl -X POST http://localhost:9080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass123"
  }'
```

5. Deploy your first function
```bash
# See api_examples.md for detailed examples
```

## Function Development

### Function Contract
Functions must follow these guidelines:
1. **Input Format**
   ```json
   {
       "direct_inputs": {
           "param1": "value1",
           "param2": "value2"
       },
       "object_inputs": {
           "file1": "function_id/object_name"
       },
       "secrets": [
           "SECRET_NAME1",
           "SECRET_NAME2"
       ]
   }
   ```

2. **Environment Variables**
   - `API_BASE_URL`: Base URL for accessing function objects
   - Secret variables: All requested secrets are injected as env vars
   - Example: A secret named "database_url" becomes "DATABASE_URL"

3. **Output Requirements**
   - Must return valid JSON via stdout
   - Errors should be returned as: `{"error": "error message"}`
   - Success response should be relevant to function purpose

4. **Object Access**
   - Objects can be accessed via HTTP GET to `$API_BASE_URL/{function_id}/{object_name}`
   - Authentication is handled automatically within the function network

5. **Execution Context**
   - Functions run in isolated containers
   - 5-minute execution timeout
   - Network access limited to internal services
   - Stateless execution model

## API Routes

### Authentication
```
POST   /auth/register           # Register new user
POST   /auth/login             # Login and get token
```

### Functions
```
GET    /api/functions          # List functions
POST   /api/functions          # Create function
GET    /api/functions/:id      # Get function details
DELETE /api/functions/:id      # Delete function
```

### Executions
```
POST   /api/executions         # Execute function
GET    /api/executions         # List executions
GET    /api/executions/:id     # Get execution status/result
```

### Function Objects
```
POST   /api/function-objects/:function_id/:name    # Upload object
GET    /api/function-objects/:function_id/:name    # Download object
DELETE /api/function-objects/:function_id/:name    # Delete object
GET    /api/function-objects/:function_id          # List objects
```

### Secrets
```
POST   /api/secrets            # Create secret
GET    /api/secrets            # List secrets
GET    /api/secrets/:id        # Get secret metadata
PUT    /api/secrets/:id        # Update secret value
DELETE /api/secrets/:id        # Delete secret
```

### Users
```
GET    /api/users              # List users (admin only)
GET    /api/users/:id          # Get user details
```

### Example Function
```python
import os
import sys
import json

def main():
    # Get input from arguments
    input_data = json.loads(sys.argv[1])
    
    # Process input
    result = {"message": "Hello " + input_data["name"]}
    
    # Return JSON output
    print(json.dumps(result))

if __name__ == "__main__":
    main()
```

## Documentation
- [API Documentation](docs/api_examples.md)
- [Function Contract](docs/function_contract.md)
- [Architecture Details](docs/architecture.md)

## License
MIT License 



## Auto Scaling

The platform includes an auto-scaling component that automatically adjusts the number of worker nodes based on NATS metrics.

### How it Works

1. **NATS Metrics Monitoring**
   - Monitors pending messages in execution queue
   - Tracks consumer count
   - Monitors message processing rates

2. **Scaling Decisions**
   - Scale Up Triggers:
     - High number of pending messages
     - Message processing delay exceeds threshold
     - Consumer to pending messages ratio is low

   - Scale Down Triggers:
     - Low queue utilization
     - Fast message processing
     - More consumers than needed for current load

3. **Implementation**
   The auto-scaler uses NATS JetStream metrics:
   ```go
   type AutoScaler struct {
       js           nats.JetStreamContext
       streamName   string
       consumerName string
   }

   func (s *AutoScaler) getMetrics() (*StreamMetrics, error) {
       // Get stream info
       stream, err := s.js.StreamInfo(s.streamName)
       if err != nil {
           return nil, err
       }

       // Get consumer info
       consumer, err := s.js.ConsumerInfo(s.streamName, s.consumerName)
       if err != nil {
           return nil, err
       }

       return &StreamMetrics{
           PendingMessages: stream.State.Msgs,
           ConsumerCount:   len(consumer.Cluster.Replicas),
           ProcessingLag:   consumer.NumPending,
       }, nil
   }
   ```

### Benefits
- No need for external metrics collection
- Real-time scaling based on actual workload
- Native integration with NATS JetStream
- Simple and efficient implementation
- Low overhead monitoring 