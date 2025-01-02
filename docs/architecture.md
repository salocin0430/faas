# FaaS Platform Documentation

## Arquitectura General

### 1. Visión General
La plataforma FaaS (Function as a Service) está construida siguiendo los principios de:
- Arquitectura Hexagonal (Ports & Adapters)
- Domain-Driven Design (DDD)
- Microservicios
- Event-Driven Architecture

### 2. Componentes Principales

#### 2.1 API Gateway (APISIX)
- Punto de entrada único
- Manejo de rutas
- Rate limiting
- Autenticación básica
- Configuración:
  ```yaml
  routes:
    - uri: /functions/*
    - uri: /executions/*
    - uri: /users/*
  ```

#### 2.2 API Service
- Manejo de requests HTTP
- Autenticación y autorización
- Gestión de funciones y ejecuciones
- Configuración:
  ```go
  type Config struct {
      ServerAddress string
      NatsURL      string
      JWTSecret    string
  }
  ```

#### 2.3 NATS JetStream
- Message broker y event store
- Streams:
  - EXECUTIONS: Cola de ejecuciones pendientes
  - RESULTS: Resultados de ejecuciones
- Configuración:
  ```go
  const (
      EXECUTIONS_SUBJECT = "executions.pending"
      WORKERS_QUEUE = "execution-workers"
  )
  ```

#### 2.4 Workers
- Procesamiento de funciones
- Escalado automático
- Gestión de contenedores Docker
- Configuración:
  ```go
  type WorkerConfig struct {
      NatsURL string
      MaxConcurrentExecutions int
  }
  ```

#### 2.5 Autoscaler
- Monitoreo de carga
- Escalado automático de workers
- Configuración:
  ```go
  type AutoscalerConfig struct {
      MinWorkers         int    // Mínimo de workers (2)
      MaxWorkers         int    // Máximo de workers (10)
      ScaleUpThreshold   int    // Mensajes por worker para escalar arriba (3)
      ScaleDownThreshold int    // Mensajes por worker para escalar abajo (2)
      CheckInterval      string // Intervalo de chequeo ("30s")
      CooldownPeriod     string // Tiempo entre escalados ("30s")
  }
  ```

### 3. Flujo de Ejecución

#### 3.1 Creación de Función
1. Cliente envía POST /functions
2. API valida y almacena metadatos
3. Función queda disponible para ejecución

#### 3.2 Ejecución de Función


### 4. Mecanismos de Seguridad

#### 4.1 Autenticación
- JWT para API
- Tokens de acceso para funciones
- Configuración:
  ```go
  type AuthConfig struct {
      JWTSecret string
      TokenExpiration time.Duration
  }
  ```

#### 4.2 Aislamiento
- Contenedores Docker para ejecución
- Límites de recursos por función
- Timeouts configurables

### 5. Escalabilidad

#### 5.1 Worker Scaling
- Escalado automático basado en carga
- Mínimo 2 workers
- Máximo 10 workers
- Métricas:
  - Mensajes pendientes
  - Workers activos
  - Mensajes por worker

#### 5.2 Queue Management
- Cola distribuida con NATS
- Rebalanceo automático
- Reintentos configurables

### 6. Monitoreo y Logging

#### 6.1 Logs
- Formato estructurado
- Niveles de log configurables
- Eventos clave registrados:
  - Inicio/fin de ejecuciones
  - Escalado de workers
  - Errores y excepciones

#### 6.2 Métricas
- Tiempo de ejecución
- Tasa de éxito/error
- Uso de recursos

### 7. Configuración

#### 7.1 Variables de Entorno
```env
# API Service
SERVER_ADDRESS=:8080
NATS_URL=nats://nats:4222
JWT_SECRET=your-secret-key

# Worker
MAX_CONCURRENT_EXECUTIONS=10

# Autoscaler
MIN_WORKERS=2
MAX_WORKERS=10
SCALE_UP_THRESHOLD=3
SCALE_DOWN_THRESHOLD=2
CHECK_INTERVAL=30s
COOLDOWN_PERIOD=30s
```

### 8. Desarrollo y Despliegue

#### 8.1 Requisitos
- Go 1.21+
- Docker
- NATS Server
- Make

#### 8.2 Comandos Make
```makefile
make build    # Construir imágenes
make run      # Ejecutar servicios
make test     # Ejecutar tests
make clean    # Limpiar recursos
```

#### 8.3 Docker Compose
- Servicios definidos
- Redes configuradas
- Volúmenes persistentes

### 9. Manejo de Errores

#### 9.1 Tipos de Errores
- Errores de validación
- Errores de ejecución
- Timeouts
- Errores de sistema

#### 9.2 Estrategias de Retry
- Backoff exponencial
- Máximo de reintentos
- Dead letter queues

### 10. Mejores Prácticas

#### 10.1 Desarrollo
- Tests unitarios
- Integración continua
- Revisión de código
- Documentación actualizada

#### 10.2 Operación
- Monitoreo proactivo
- Backups regulares
- Actualizaciones planificadas
- Gestión de incidentes 

### 11. NATS JetStream Configuración Detallada

#### 11.1 Buckets y Streams
```go
// Buckets
const (
    EXECUTIONS_BUCKET = "executions"
    FUNCTIONS_BUCKET  = "functions"
    USERS_BUCKET      = "users"
)

// Streams
const (
    EXECUTIONS_STREAM  = "EXECUTIONS"
    EXECUTIONS_SUBJECT = "executions.pending"
    RESULTS_SUBJECT    = "executions.results"
    WORKERS_QUEUE      = "execution-workers"
)

// Configuración del Stream
stream, err := js.AddStream(&nats.StreamConfig{
    Name:     EXECUTIONS_STREAM,
    Subjects: []string{EXECUTIONS_SUBJECT, RESULTS_SUBJECT},
    Storage:  nats.FileStorage,
    MaxAge:   24 * time.Hour,  // Retención de mensajes
    Replicas: 1,              // Para desarrollo
})
```

#### 11.2 Worker Configuración
```go
// Timeouts y Límites
const (
    EXECUTION_TIMEOUT = 5 * time.Minute  // Timeout por función
    MAX_RETRIES = 3                      // Reintentos por ejecución
    RETRY_DELAY = 5 * time.Second        // Delay entre reintentos
)

// Configuración de Consumer
sub, err := js.QueueSubscribe(
    EXECUTIONS_SUBJECT,
    WORKERS_QUEUE,
    handler,
    nats.ManualAck(),                // Ack manual
    nats.AckWait(1 * time.Minute),   // Tiempo de espera para ack
    nats.MaxDeliver(3),              // Máximo de reenvíos
)
```

### 12. API Endpoints

#### 12.1 Gestión de Usuarios
```http
POST   /api/v1/users/register     # Registro de usuario
POST   /api/v1/users/login        # Login de usuario
GET    /api/v1/users/me          # Información del usuario actual
PUT    /api/v1/users/me          # Actualizar usuario
DELETE /api/v1/users/me          # Eliminar usuario
```

#### 12.2 Gestión de Funciones
```http
POST   /api/v1/functions         # Crear función
GET    /api/v1/functions        # Listar funciones
GET    /api/v1/functions/{id}   # Obtener función
PUT    /api/v1/functions/{id}   # Actualizar función
DELETE /api/v1/functions/{id}   # Eliminar función
```

#### 12.3 Ejecuciones
```http
POST   /api/v1/executions                # Ejecutar función
GET    /api/v1/executions               # Listar ejecuciones
GET    /api/v1/executions/{id}          # Estado de ejecución
DELETE /api/v1/executions/{id}          # Cancelar ejecución
GET    /api/v1/executions/{id}/logs     # Logs de ejecución
```

#### 12.4 Monitoreo y Administración
```http
GET    /api/v1/health           # Estado del servicio
GET    /api/v1/metrics         # Métricas del sistema
GET    /api/v1/workers/status  # Estado de workers
```

### 13. Configuración de Contenedores

#### 13.1 Límites de Recursos
```go
// Configuración por defecto de contenedores
hostConfig := &container.HostConfig{
    Resources: container.Resources{
        Memory:    512 * 1024 * 1024,  // 512MB RAM
        NanoCPUs:  1000000000,         // 1 CPU
        PidsLimit: &pidsLimit,         // Límite de procesos
    },
    NetworkMode: "none",              // Sin acceso a red
    AutoRemove:  true,               // Eliminar al terminar
}
```

#### 13.2 Timeouts y Límites
```go
const (
    CONTAINER_STARTUP_TIMEOUT = 30 * time.Second
    EXECUTION_TIMEOUT = 5 * time.Minute
    MAX_OUTPUT_SIZE = 1024 * 1024  // 1MB
)
```

### 14. Estructura de Datos

#### 14.1 Función
```go
type Function struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    ImageURL    string    `json:"image_url"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Owner       string    `json:"owner"`
    Description string    `json:"description"`
    Timeout     int      `json:"timeout"`
}
```

#### 14.2 Ejecución
```go
type Execution struct {
    ID         string    `json:"id"`
    FunctionID string    `json:"function_id"`
    Status     string    `json:"status"`
    Input      string    `json:"input"`
    Output     string    `json:"output"`
    Error      string    `json:"error"`
    StartTime  time.Time `json:"start_time"`
    EndTime    time.Time `json:"end_time"`
    Duration   float64   `json:"duration"`
}
```

### 15. Códigos de Estado y Errores

#### 15.1 Estados de Ejecución
```go
const (
    StatusPending   = "pending"
    StatusRunning   = "running"
    StatusComplete  = "complete"
    StatusFailed    = "failed"
    StatusCancelled = "cancelled"
    StatusTimeout   = "timeout"
)
```

#### 15.2 Códigos de Error
```go
const (
    ErrFunctionNotFound    = "function_not_found"
    ErrExecutionNotFound   = "execution_not_found"
    ErrInvalidInput        = "invalid_input"
    ErrExecutionTimeout    = "execution_timeout"
    ErrInternalError       = "internal_error"
    ErrUnauthorized        = "unauthorized"
) 
```

## Object Storage y Comunicación entre Funciones

### Object Storage
El sistema incluye un mecanismo de almacenamiento de objetos que permite:
- Almacenar archivos asociados a funciones
- Acceder a estos archivos desde las funciones durante la ejecución
- Gestionar el ciclo de vida de los objetos

#### Componentes Clave:
1. **Object Repository**: 
   - Almacena objetos usando NATS KV
   - Maneja metadata y contenido binario
   - Organiza objetos por función

2. **Object Service**:
   - Gestiona operaciones CRUD de objetos
   - Valida permisos y tipos de archivo
   - Maneja límites y cuotas

3. **API Endpoints**:
   - `/api/function-objects/:function_id/:name` para operaciones CRUD
   - Soporta subida y descarga de archivos
   - Autenticación mediante JWT

### Comunicación Interna
Las funciones pueden acceder a los objetos a través de una red interna:

1. **Red Docker**:
   - Nombre: "apisix"
   - Conecta todos los servicios
   - Aislamiento y seguridad

2. **Acceso a Objetos**:
   - URL base interna: `http://api:9080/api`
   - No requiere autenticación en red interna
   - Configurado vía variables de entorno

### Ejemplo: PDF Processor

Este ejemplo demuestra la integración completa del sistema:

1. **Componentes**:


2. **Flujo de Datos**:
   - Cliente sube PDF vía API Gateway
   - Objeto almacenado en NATS KV
   - Función accede al PDF vía red interna
   - Resultado devuelto al cliente

3. **Características**:
   - Procesamiento asíncrono
   - Almacenamiento persistente
   - Comunicación segura
   - Escalabilidad horizontal

### Consideraciones de Diseño

1. **Seguridad**:
   - Autenticación externa vía JWT
   - Red interna confiable
   - Aislamiento de contenedores

2. **Rendimiento**:
   - Acceso directo a objetos
   - Caché de objetos (futuro)
   - Optimización de red

3. **Escalabilidad**:
   - Múltiples workers
   - Distribución de carga
   - Replicación de datos 