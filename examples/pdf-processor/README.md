# PDF Processor Function

Esta función procesa archivos PDF y extrae texto de las páginas especificadas.

## Construcción

```bash
docker build -t salocin0430/pdf-processor:v1 .
docker push salocin0430/pdf-processor:v1
```

## Uso

1. Crear la función:
```bash
curl -X POST http://localhost:9080/api/functions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "pdf-processor",
    "image_url": "salocin0430/pdf-processor:v1",
    "description": "Procesa archivos PDF y extrae texto"
  }'
```

2. Subir un PDF:
```bash
curl -X POST "http://localhost:9080/api/function-objects/$FUNCTION_ID/document.pdf" \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@./document.pdf"
```

3. Ejecutar:
```bash
curl -X POST http://localhost:9080/api/executions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "function_id": "'$FUNCTION_ID'",
    "input": {
        "direct_inputs": {
            "pages": 2
        },
        "object_inputs": {
            "pdf_file": "function-objects/'$FUNCTION_ID'/document.pdf"
        }
    }
  }'
``` 