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
    "input": "Hello World"
  }'

# Respuesta esperada (inmediata)
{
    "id": "98765432-abcd-efgh-ijkl-123456789000",
    "function_id": "123e4567-e89b-12d3-a456-426614174000",
    "status": "pending",
    "input": "Hello World",
    "created_at": "2024-03-21T10:30:00Z"
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
curl -X POST http://localhost:8080/api/executions \
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