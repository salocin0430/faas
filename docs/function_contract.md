# FaaS Function Contract

This document describes the contract that all functions must follow to be compatible with our FaaS system.

## 1. General Rules

### Input and argv (Command Line Arguments)
- Function must accept a single string argument in JSON format
- The argument is passed as a command line argument, known as "argv"
- The input must be an escaped JSON string when sending to API
- Input JSON structure:
  ```json
  {
      "direct_inputs": {
          "param1": "value1",
          "param2": 123,
          "param3": true,
          "param4": {
              "nested1": "value",
              "nested2": 456
          },
          "param5": ["a", "b", "c"]
      },
      "object_inputs": {
          "file1": "function_id/document.pdf",
          "config": "function_id/config.json"
      },
      "secrets": [
          "API_KEY",
          "DATABASE_URL"
      ]
  }
  ```

### Important Notes About Input
1. The entire input JSON must be escaped as a string when sending to API
2. All sections (direct_inputs, object_inputs, secrets) are optional
3. Object paths must follow format: `function_id/object_name`
4. Secret names are converted to environment variables

### Output Format
- Result must be written to stdout
- Must be a valid JSON string
- Must not contain logs or additional messages
- Two possible formats:

Success response:
```json
{
    "result": "any valid json value",
    "metadata": {
        "duration_ms": 1500,
        "items_processed": 10
    }
}
```

Error response:
```json
{
    "error": "detailed error message"
}
```

### Logs and Errors
- All logs must be written to stderr
- Errors must be reported to stderr
- Exit code must be 0 for success, non-zero for error

### Example API Usage
```bash
# Wrong - Input not escaped (don't do this)
curl -X POST http://localhost:9080/api/executions \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "function_id": "123abc",
    "input": {
        "direct_inputs": {"name": "John"},
        "secrets": ["API_KEY"]
    }
  }'

# Correct - Input properly escaped (do this)
curl -X POST http://localhost:9080/api/executions \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "function_id": "123abc",
    "input": "{\"direct_inputs\":{\"name\":\"John\"},\"secrets\":[\"API_KEY\"]}"
  }'
```

## 2. Object Access

### Environment Variables
Functions have access to:
```bash
API_BASE_URL="http://api:8080/api/function-objects"  # Base URL for object access
```

### Accessing Objects
Objects can be accessed via HTTP GET requests to the internal API:
```python
url = f"{os.getenv('API_BASE_URL')}/{object_ref}"
response = requests.get(url)  # No authentication needed in internal network
```

## 3. Examples by Language

### Python
```python
#!/usr/bin/env python3
import os
import sys
import json
import requests
import logging

# Configure logging to stderr
logging.basicConfig(
    stream=sys.stderr,
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

def process_input(data):
    try:
        # Get object if specified
        if "object_inputs" in data:
            file_ref = data["object_inputs"]["file1"]
            url = f"{os.getenv('API_BASE_URL')}/{file_ref}"
            response = requests.get(url)
            if response.status_code != 200:
                raise Exception(f"Error getting file: {response.status_code}")
            # Process file content
            content = response.content
            # ... process content ...

        # Process direct inputs
        result = data["direct_inputs"]["value"] * 2
        return {"result": result}
    except Exception as e:
        return {"error": str(e)}

def main():
    try:
        # Verify arguments
        if len(sys.argv) != 2:
            raise ValueError("Exactly one JSON argument required")

        # Get JSON from command line argument
        json_input = sys.argv[1]
        logging.info(f"Received input: {json_input}")

        # Parse JSON input
        input_data = json.loads(json_input)
        
        # Process
        result = process_input(input_data)

        # Write result to stdout
        print(json.dumps(result))
        sys.exit(0)

    except Exception as e:
        logging.error(f"Error: {str(e)}")
        print(json.dumps({"error": str(e)}))
        sys.exit(1)

if __name__ == "__main__":
    main()
```

### Node.js
```javascript
#!/usr/bin/env node
const process = require('process');
const axios = require('axios');

// Configure logging
function log(message) {
    console.error(`[${new Date().toISOString()}] ${message}`);
}

// Main processing function
async function processInput(data) {
    try {
        // Get object if specified
        if (data.object_inputs) {
            const fileRef = data.object_inputs.file1;
            const url = `${process.env.API_BASE_URL}/${fileRef}`;
            
            const response = await axios.get(url, {
                responseType: 'arraybuffer'
            });
            
            if (response.status !== 200) {
                throw new Error(`Error getting file: ${response.status}`);
            }
            
            // Process file content
            const content = response.data;
            // ... process content ...
        }

        // Process direct inputs
        const result = data.direct_inputs.value * 2;
        return { result };
    } catch (error) {
        return { error: error.message };
    }
}

async function main() {
    try {
        // Verify arguments (remember Node.js has 2 extra args)
        if (process.argv.length !== 3) {
            throw new Error('Exactly one JSON argument required');
        }

        // Get JSON from command line argument
        const jsonInput = process.argv[2];
        log(`Received input: ${jsonInput}`);

        // Parse JSON input
        const inputData = JSON.parse(jsonInput);
        
        // Process
        const result = await processInput(inputData);

        // Result to stdout
        console.log(JSON.stringify(result));
        process.exit(0);

    } catch (error) {
        log(`Error: ${error.message}`);
        console.log(JSON.stringify({ error: error.message }));
        process.exit(1);
    }
}

main();
```

### Go
```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
)

// Input represents the expected JSON input structure
type Input struct {
    DirectInputs struct {
        Value int `json:"value"`
    } `json:"direct_inputs"`
    ObjectInputs struct {
        File1 string `json:"file1"`
    } `json:"object_inputs"`
}

// Output represents the JSON output structure
type Output struct {
    Result interface{} `json:"result,omitempty"`
    Error  string     `json:"error,omitempty"`
}

func processInput(input Input) Output {
    // Get object if specified
    if input.ObjectInputs.File1 != "" {
        apiBaseURL := os.Getenv("API_BASE_URL")
        url := fmt.Sprintf("%s/%s", apiBaseURL, input.ObjectInputs.File1)

        resp, err := http.Get(url)
        if err != nil {
            return Output{Error: fmt.Sprintf("Error getting file: %v", err)}
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            return Output{Error: fmt.Sprintf("Error getting file: %d", resp.StatusCode)}
        }

        // Read and process file content
        content, err := io.ReadAll(resp.Body)
        if err != nil {
            return Output{Error: fmt.Sprintf("Error reading content: %v", err)}
        }
        // ... process content ...
    }

    // Process direct inputs
    result := input.DirectInputs.Value * 2
    return Output{Result: result}
}

func main() {
    // Configure log for stderr
    log.SetOutput(os.Stderr)
    log.SetFlags(log.LstdFlags | log.Lmicroseconds)

    // Verify arguments
    if len(os.Args) != 2 {
        fmt.Println(json.Marshal(Output{Error: "Exactly one JSON argument required"}))
        os.Exit(1)
    }

    // Get JSON from command line argument
    jsonInput := os.Args[1]
    log.Printf("Received input: %s", jsonInput)

    // Parse JSON input
    var input Input
    if err := json.Unmarshal([]byte(jsonInput), &input); err != nil {
        log.Printf("Error parsing JSON: %v", err)
        json.NewEncoder(os.Stdout).Encode(Output{Error: err.Error()})
        os.Exit(1)
    }

    // Process
    result := processInput(input)

    // Write result to stdout
    if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
        log.Printf("Error encoding result: %v", err)
        os.Exit(1)
    }
}
```
## 4. Runtime Environment

### Network
- Functions run in the "apisix" Docker network
- Direct access to internal API (no auth needed)
- Environment variables for configuration

### Resource Limits
- Execution timeout: 5 minutes
- Maximum output size: 1MB 
- Memory: 512MB (default) TODO
- CPU: 1 core (default)  TODO

## 5. Best Practices

### Error Handling
- Always use try-catch blocks
- Return structured error messages
- Log relevant information

### Resource Management
- Close files and connections
- Clean up temporary resources
- Handle memory efficiently

### Security
- Validate all inputs
- Don't trust external data TODO
- Handle errors appropriately
