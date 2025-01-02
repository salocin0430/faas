# FaaS API Usage Examples

This document provides examples of how to interact with the API using cURL.

## 1. User Management

### Create User
```bash
# Registro de usuario
curl -X POST http://localhost:9080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser5",
    "password": "testpass123",
    "role": "user"
  }'

# Successful Response
{
    "id": "user123",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2023-11-22T10:30:00Z"
}
```

### Login
```bash
# Login
curl -X POST http://localhost:9080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser5",
    "password": "testpass123"
  }'

# Successful Response
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzQ5NjI3MjUsInJvbGUiOiJ1c2VyIiwic3ViIjoiOGY2ZTVhNDMtMDdjYi00YmVkLWI2Y2MtYjBhNmM3YmRhOTI3IiwidXNlciI6InRlc3R1c2VyNSIsInVzZXIta2V5IjoidGVzdHVzZXI1In0.pi7Z73vAprLCxycahgAYKBI7sgcPMbVzAmLi9IL0-Bs"
}
```
## 1.2 List users
# Listar usuarios (ruta protegida)
```bash
curl -X GET http://localhost:9080/api/users \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzQ5NzkzMTEsImtleSI6ImZhYXNhcHAta2V5Iiwicm9sZSI6InVzZXIiLCJzdWIiOiI4ZjZlNWE0My0wN2NiLTRiZWQtYjZjYy1iMGE2YzdiZGE5MjciLCJ1c2VyIjoidGVzdHVzZXI1In0.J__dllWYcR214AYoFRNV8bYNt47s_afLvOgSjvZOfsE"
  
```
## 2. Function Management

### Create Function
```bash
# Store token for subsequent requests
TOKEN="eyJhbGciOiJIUzI1NiIs..."

# 2. Crear función usando el token
curl -X POST http://localhost:9080/api/functions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzQ5NzA3MjUsInJvbGUiOiJ1c2VyIiwic3ViIjoiOGY2ZTVhNDMtMDdjYi00YmVkLWI2Y2MtYjBhNmM3YmRhOTI3IiwidXNlciI6InRlc3R1c2VyNSJ9.cQMFaIWxkWNdQzGiBA0blE9SMTT-WIG4xjZeEtJAjRI" \
  -d '{
    "name": "hello-world",
    "image_url": "docker.io/library/hello-world:latest",
    "description": "A simple hello world function"
  }'

# Successful Response
{
    "id": "func123",
    "name": "multiply",
    "image_url": "docker.io/myrepo/multiply:latest",
    "description": "Multiplies a number by 2",
    "created_at": "2023-11-22T10:35:00Z"
}
```

### List Functions
```bash
curl -X GET http://localhost:8080/api/functions \
  -H "Authorization: Bearer $TOKEN"

# Successful Response
{
    "functions": [
        {
            "id": "func123",
            "name": "multiply",
            "image_url": "docker.io/myrepo/multiply:latest",
            "description": "Multiplies a number by 2",
            "created_at": "2023-11-22T10:35:00Z"
        }
    ]
}
```

### Get Specific Function
```bash
curl -X GET http://localhost:8080/api/functions/func123 \
  -H "Authorization: Bearer $TOKEN"

# Successful Response
{
    "id": "func123",
    "name": "multiply",
    "image_url": "docker.io/myrepo/multiply:latest",
    "description": "Multiplies a number by 2",
    "created_at": "2023-11-22T10:35:00Z"
}
```

## 3. Function Execution

### Execute Function
```bash
curl -X POST http://localhost:9080/api/executions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <tu-token-jwt>" \
  -d '{
    "function_id": "123e4567-e89b-12d3-a456-426614174000",
    "input": "{\"direct_inputs\": {\"name\": \"John\", \"age\": 30}, \"object_inputs\": {\"file\": \"123e4567/document.pdf\"}, \"secrets\": [\"API_KEY\", \"DATABASE_URL\"]}"
  }'

# Respuesta esperada (inmediata)
{
    "id": "98765432-abcd-efgh-ijkl-123456789000",
    "function_id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "pending",
    "input": "{\"direct_inputs\":{\"name\":\"John\",\"age\":30},\"object_inputs\":{\"file\":\"123e4567/document.pdf\"},\"secrets\":[\"API_KEY\",\"DATABASE_URL\"]}",
    "created_at": "2024-03-21T10:30:00Z"
}

# Formato del input (debe ser un string JSON escapado):
{
    "direct_inputs": {          // Parámetros directos para la función
        "param1": "value1",
        "param2": value2
    },
    "object_inputs": {          // Referencias a archivos almacenados
        "file_alias": "function_id/object_name"
    },
    "secrets": [               // Lista de secrets a inyectar como variables de entorno
        "SECRET_NAME1",
        "SECRET_NAME2"
    ]
}
```

### Check Execution Status
```bash
curl -X GET http://localhost:9080/api/executions/98765432-abcd-efgh-ijkl-123456789000 \
  -H "Authorization: Bearer <tu-token-jwt>"

# Successful Response
{
    "id": "exec123",
    "function_id": "func123",
    "status": "completed",
    "input": "{\"value\": 21}",
    "output": "{\"result\": 42}",
    "created_at": "2023-11-22T10:40:00Z",
    "completed_at": "2023-11-22T10:40:02Z"
}
```

### List Executions
```bash
curl -X GET http://localhost:8080/api/executions \
  -H "Authorization: Bearer $TOKEN"

# Successful Response
{
    "executions": [
        {
            "id": "exec123",
            "function_id": "func123",
            "status": "completed",
            "created_at": "2023-11-22T10:40:00Z",
            "completed_at": "2023-11-22T10:40:02Z"
        }
    ]
}
```

## 4. Complete Flow Example

### Create and Execute a Function
```bash
# 1. Login
TOKEN=$(curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "secretpass123"
  }' | jq -r '.token')

# 2. Create Function
FUNCTION_ID=$(curl -X POST http://localhost:8080/api/functions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "multiply",
    "image_url": "docker.io/myrepo/multiply:latest",
    "description": "Multiplies a number by 2"
  }' | jq -r '.id')

# 3. Execute Function
EXECUTION_ID=$(curl -X POST http://localhost:8080/api/executions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"function_id\": \"$FUNCTION_ID\",
    \"input\": \"{\\\"value\\\": 21}\"
  }" | jq -r '.execution_id')

# 4. Wait and Get Result
sleep 2
curl -X GET "http://localhost:8080/api/executions/$EXECUTION_ID" \
  -H "Authorization: Bearer $TOKEN"
```

## 5. Error Handling

### Authentication Error Example
```bash
curl -X GET http://localhost:8080/api/functions \
  -H "Authorization: Bearer invalid_token"

# Response
{
    "error": "invalid_token",
    "message": "Invalid or expired token"
}
```

### Execution Error Example
```bash
curl -X POST http://localhost:9080/api/executions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "function_id": "func123",
    "input": "invalid json"
  }'

# Response
{
    "error": "invalid_input",
    "message": "Input must be valid JSON"
}
```

## 6. Important Notes

1. **Required Headers**:
   - `Content-Type: application/json` for requests with body
   - `Authorization: Bearer <token>` for authenticated endpoints

2. **Authentication**:
   - Get token through login
   - Include token in all subsequent requests
   - Tokens expire after a certain time

3. **Asynchronous Execution**:
   - Function execution is asynchronous
   - Returns execution ID immediately
   - Use status endpoint to check completion

## 4. Function Objects Management

### Upload Object
```bash
# Subir un archivo a una función
curl -X POST http://localhost:9080/api/function-objects/123/objects/dataset-1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: multipart/form-data" \
  -F "file=@/path/to/your/file.csv"

# Successful Response
{
    "id": "obj_123",
    "name": "dataset-1",
    "size": 1048576,
    "content_type": "text/csv",
    "created_at": "2024-03-21T15:30:00Z"
}
```

### List Objects
```bash
# Listar objetos de una función
curl -X GET http://localhost:9080/api/function-objects/123/objects \
  -H "Authorization: Bearer $TOKEN"

# Successful Response
{
    "objects": [
        {
            "id": "obj_123",
            "name": "dataset-1",
            "size": 1048576,
            "content_type": "text/csv",
            "created_at": "2024-03-21T15:30:00Z"
        }
    ]
}
```

### Get Object
```bash
# Descargar un objeto
curl -X GET http://localhost:9080/api/function-objects/123/objects/dataset-1 \
  -H "Authorization: Bearer $TOKEN" \
  -O dataset-1.csv
```

### Delete Object
```bash
# Eliminar un objeto
curl -X DELETE http://localhost:9080/api/function-objects/123/objects/dataset-1 \
  -H "Authorization: Bearer $TOKEN"
```

### Example Usage with Function Execution
```bash
# 1. Subir dataset
curl -X POST http://localhost:9080/api/function-objects/123/objects/training-data \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/dataset.csv"

# 2. Ejecutar función usando el objeto
curl -X POST http://localhost:9080/api/executions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "function_id": "123",
    "input": {
        "direct_inputs": {
            "model": "regression",
            "iterations": 100
        },
        "object_inputs": {
            "dataset": "function-objects/123/objects/training-data"
        }
    }
  }'
```

## 8. PDF Processor Example

Este ejemplo muestra cómo crear y usar una función que procesa archivos PDF.

### 1. Create PDF Processor Function
```bash
curl -X POST http://localhost:9080/api/functions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "pdf-processor",
    "image_url": "salocin0430/pdf-processor:v2",
    "description": "Process PDF files and extract text"
  }'

# Response
{
    "id": "15f2a713-135e-422c-90ae-57779fe7b4ad",
    "name": "pdf-processor",
    "image_url": "salocin0430/pdf-processor:v2",
    "description": "Process PDF files and extract text",
    "created_at": "2024-01-02T15:30:00Z"
}
```

### 2. Subir archivo PDF
```bash
curl -X POST "http://localhost:9080/api/function-objects/15f2a713-135e-422c-90ae-57779fe7b4ad/PRINCE2" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@./PRINCE2.pdf"

# Response
{
    "id": "obj_123",
    "name": "PRINCE2",
    "size": 2048576,
    "content_type": "application/pdf",
    "created_at": "2024-01-02T15:31:00Z"
}
```

### 3. Ejecutar procesamiento
```bash
curl -X POST http://localhost:9080/api/executions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "function_id": "15f2a713-135e-422c-90ae-57779fe7b4ad",
    "input": "{\"direct_inputs\": {\"pages\": 2}, \"object_inputs\": {\"pdf_file\": \"15f2a713-135e-422c-90ae-57779fe7b4ad/PRINCE2\"}}"
  }'

# Response
{
    "id": "exec_456",
    "function_id": "15f2a713-135e-422c-90ae-57779fe7b4ad",
    "status": "pending",
    "created_at": "2024-01-02T15:32:00Z"
}
```

### 4. Verificar resultado
```bash
curl -X GET "http://localhost:9080/api/executions/exec_456" \
  -H "Authorization: Bearer $TOKEN"

# Response
{
    "id": "exec_456",
    "function_id": "15f2a713-135e-422c-90ae-57779fe7b4ad",
    "status": "completed",
    "input": "{\"direct_inputs\":{\"pages\":2},\"object_inputs\":{\"pdf_file\":\"15f2a713-135e-422c-90ae-57779fe7b4ad/PRINCE2\"}}",
    "output": "{\"pages_processed\":2,\"total_pages\":245,\"text\":\"... primeros 1000 caracteres del PDF ...\"}",
    "created_at": "2024-01-02T15:32:00Z",
    "completed_at": "2024-01-02T15:32:05Z"
}
```

### Notas importantes:
1. La función espera un PDF como input
2. El parámetro `pages` determina cuántas páginas procesar
3. El texto extraído se limita a los primeros 1000 caracteres
4. La función usa la red interna "apisix" para comunicarse con el API
5. No es necesario pasar tokens internamente ya que está en la misma red

## 9. Secrets Management

### Create Secret
```bash
curl -X POST http://localhost:9080/api/secrets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "DATABASE_URL",
    "value": "postgresql://user:pass@localhost:5432/db"
  }'

# Response
{
    "id": "sec_abc123",
    "name": "DATABASE_URL",
    "created_at": "2024-01-03T10:00:00Z",
    "updated_at": "2024-01-03T10:00:00Z"
}
```

### List Secrets
```bash
curl -X GET http://localhost:9080/api/secrets \
  -H "Authorization: Bearer $TOKEN"

# Response
{
    "secrets": [
        {
            "id": "sec_abc123",
            "name": "DATABASE_URL",
            "created_at": "2024-01-03T10:00:00Z",
            "updated_at": "2024-01-03T10:00:00Z"
        },
        {
            "id": "sec_def456",
            "name": "API_KEY",
            "created_at": "2024-01-03T11:00:00Z",
            "updated_at": "2024-01-03T11:00:00Z"
        }
    ]
}
```

### Get Secret
```bash
curl -X GET http://localhost:9080/api/secrets/sec_abc123 \
  -H "Authorization: Bearer $TOKEN"

# Response
{
    "id": "sec_abc123",
    "name": "DATABASE_URL",
    "created_at": "2024-01-03T10:00:00Z",
    "updated_at": "2024-01-03T10:00:00Z"
}
```

### Update Secret
```bash
curl -X PUT http://localhost:9080/api/secrets/sec_abc123 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "value": "postgresql://newuser:newpass@localhost:5432/newdb"
  }'

# Response
{
    "id": "sec_abc123",
    "name": "DATABASE_URL",
    "created_at": "2024-01-03T10:00:00Z",
    "updated_at": "2024-01-03T12:00:00Z"
}
```

### Delete Secret
```bash
curl -X DELETE http://localhost:9080/api/secrets/sec_abc123 \
  -H "Authorization: Bearer $TOKEN"

# Response: 204 No Content
```

### Notes:
1. All secret operations require authentication
2. Secret values are only shown during creation and update
3. List and Get operations only return metadata (no values)
4. Secrets are scoped to the authenticated user
5. Secret names must be unique per user
```

## 10. Secret Printer Example

This example demonstrates how to use secrets and make HTTP requests in a function.

### 1. Create Function
```bash
curl -X POST http://localhost:9080/api/functions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "secret-printer",
    "image_url": "salocin0430/secret-printer:v1",
    "description": "Prints secrets and makes HTTP requests"
  }'
```

### 2. Create Some Secrets
```bash
# Create API key secret
curl -X POST http://localhost:9080/api/secrets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "API_KEY",
    "value": "my-super-secret-key"
  }'

# Create another secret
curl -X POST http://localhost:9080/api/secrets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "SERVICE_URL",
    "value": "https://api.internal.service"
  }'
```

### 3. Execute Function
```bash
curl -X POST http://localhost:9080/api/executions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "function_id": "func123",
    "input": {
        "direct_inputs": {
            "api_url": "https://jsonplaceholder.typicode.com/todos/1"
        },
        "secrets": [
            "API_KEY",
            "SERVICE_URL"
        ]
    }
  }'
```

### 4. Expected Output
```json
{
    "secrets_found": 2,
    "api_status": 200,
    "stdout": "=== Secrets ===\nAPI_KEY: my-super-secret-key\nSERVICE_URL: https://api.internal.service\n\n=== API Response ===\n{\n  \"userId\": 1,\n  \"id\": 1,\n  \"title\": \"delectus aut autem\",\n  \"completed\": false\n}"
}
```