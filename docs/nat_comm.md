# NATS Commands Guide

## Key-Value Operations

### List and Inspect KV Stores
```bash
# Listar todos los KV stores
nats kv ls

# Ver detalles de un KV store específico
nats kv info functions
nats kv info executions
nats kv info users
```

### View KV Contents
```bash
# Ver todas las keys en un bucket
nats kv ls functions
nats kv ls executions
nats kv ls users

# Ver valor de una key específica
nats kv get functions <key-id>
nats kv get executions <key-id>
nats kv get users <key-id>

# Observar cambios en tiempo real
nats kv watch functions
nats kv watch executions
nats kv watch users
```

## Stream Operations

### List and Inspect Streams
```bash
# Listar todos los streams
nats stream ls

# Ver detalles del stream de ejecuciones
nats stream info EXECUTIONS

# Ver mensajes en el stream
nats stream view EXECUTIONS

# Ver mensajes en tiempo real (sin consumir)
nats sub "executions.>"
nats sub "executions.pending"
```

### Consumer Operations
```bash
# Listar consumers de un stream
nats consumer ls EXECUTIONS

# Crear un consumer
nats consumer add EXECUTIONS test-consumer --filter executions.pending --ack none

# Leer mensajes (sin consumir)
nats consumer next EXECUTIONS test-consumer --no-ack

# Ver mensajes sin consumirlos
nats consumer peek EXECUTIONS test-consumer
```

## Monitoring Commands
```bash
# Ver estadísticas generales
nats server report

# Ver información de conexiones
nats connection list

# Monitorear latencia
nats latency --server nats://localhost:4222
```

## Debug Tips
1. Usar `watch` para cambios en KV
2. Usar `sub` para mensajes nuevos
3. Usar `consumer peek` para ver mensajes sin consumirlos
4. Usar `view` para ver histórico completo
