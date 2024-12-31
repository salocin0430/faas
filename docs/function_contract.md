# FaaS Function Contract

This document describes the contract that all functions must follow to be compatible with our FaaS system.

## 1. General Rules

### Input and argv (Command Line Arguments)
- Function must accept a single string argument in JSON format
- The argument is passed as a command line argument, known as "argv":
  - In most languages, argv is an array/list of strings containing program arguments
  - argv[0] is typically the program name/path
  - argv[1] contains our JSON input string
  - Example: `./myfunction '{"value": 42}'`

### How argv Works in Different Languages

#### Python
```python
import sys

# sys.argv is a list of strings
# sys.argv[0] = program name (e.g., "./myfunction")
# sys.argv[1] = JSON input string
json_input = sys.argv[1]  # e.g., '{"value": 42}'
```

#### Node.js
```javascript
// process.argv is an array of strings
// process.argv[0] = node executable path
// process.argv[1] = script path
// process.argv[2] = JSON input string (note the index is 2, not 1)
const json_input = process.argv[2]  // e.g., '{"value": 42}'
```

#### Go
```go
// os.Args is a slice of strings
// os.Args[0] = program name
// os.Args[1] = JSON input string
json_input := os.Args[1]  // e.g., '{"value": 42}'
```

### Output
- Result must be written to stdout
- Must be a valid JSON string
- Must not contain logs or additional messages

### Logs and Errors
- All logs must be written to stderr
- Errors must be reported to stderr
- Exit code must be 0 for success, non-zero for error

## 2. Examples by Language

### Python
```python
#!/usr/bin/env python3
import sys
import json
import logging

# Configure logging to stderr
logging.basicConfig(
    stream=sys.stderr,
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

def process_input(data):
    # Your logic here
    return {"result": data["value"] * 2}

def main():
    try:
        # Verify arguments
        if len(sys.argv) != 2:
            raise ValueError("Exactly one JSON argument required")

        # Get JSON from command line argument
        json_input = sys.argv[1]
        logging.info(f"Received argv[1]: {json_input}")

        # Parse JSON input
        input_data = json.loads(json_input)
        logging.info(f"Parsed input: {input_data}")

        # Process
        result = process_input(input_data)

        # Write result to stdout
        print(x`.dumps(result))
        sys.exit(0)

    except Exception as e:
        # Error logs to stderr
        logging.error(f"Error: {str(e)}")
        # Error result to stdout
        print(json.dumps({"error": str(e)}))
        sys.exit(1)

if __name__ == "__main__":
    main()
```

### Node.js
```javascript
#!/usr/bin/env node
const process = require('process');

// Main processing function
function processInput(data) {
    // Your logic here
    return { result: data.value * 2 };
}

// Logs go to stderr
function log(message) {
    console.error(`[${new Date().toISOString()}] ${message}`);
}

try {
    // Verify arguments (remember Node.js has 2 extra args)
    if (process.argv.length !== 3) {
        throw new Error('Exactly one JSON argument required');
    }

    // Get JSON from command line argument
    const jsonInput = process.argv[2];  // Note: using index 2
    log(`Received argv[2]: ${jsonInput}`);

    // Parse JSON input
    const inputData = JSON.parse(jsonInput);
    log(`Parsed input: ${JSON.stringify(inputData)}`);

    // Process
    const result = processInput(inputData);

    // Result to stdout
    console.log(JSON.stringify(result));
    process.exit(0);

} catch (error) {
    // Error to stderr
    log(`Error: ${error.message}`);
    // Error result to stdout
    console.log(JSON.stringify({ error: error.message }));
    process.exit(1);
}
```

### Go
```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
)

// Input represents the JSON input structure
type Input struct {
    Value int `json:"value"`
}

// Output represents the JSON output structure
type Output struct {
    Result int    `json:"result,omitempty"`
    Error  string `json:"error,omitempty"`
}

func processInput(input Input) Output {
    // Your logic here
    return Output{
        Result: input.Value * 2,
    }
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
    log.Printf("Received argv[1]: %s", jsonInput)

    // Parse JSON input
    var input Input
    if err := json.Unmarshal([]byte(jsonInput), &input); err != nil {
        log.Printf("Error parsing JSON: %v", err)
        json.NewEncoder(os.Stdout).Encode(Output{Error: err.Error()})
        os.Exit(1)
    }

    // Log parsed input
    log.Printf("Parsed input: %+v", input)

    // Process
    result := processInput(input)

    // Write result to stdout
    json.NewEncoder(os.Stdout).Encode(result)
}
```

## 3. Usage Examples

### Command Line Execution
```bash
# Python function
python3 function.py '{"value": 21}'

# Node.js function
node function.js '{"value": 21}'

# Go function (after compilation)
./function '{"value": 21}'

# Expected output in all cases:
{"result": 42}
```

### Successful Call
```bash
# Input
./function '{"value": 21}'

# Stdout (result)
{"result": 42}

# Stderr (logs)
2023-11-22 10:30:15 - INFO - Received argv[1]: {"value": 21}
2023-11-22 10:30:15 - INFO - Parsed input: {"value": 21}
```

### Error Call
```bash
# Input
./function 'invalid json'

# Stdout (result)
{"error": "Invalid JSON input"}

# Stderr (logs)
2023-11-22 10:30:20 - ERROR - Error parsing JSON: invalid character 'i' looking for beginning of value
```

## 4. Important Considerations

1. **Command Line Arguments (argv)**
   - Always check argument count before accessing
   - Remember Node.js has different argv indexing
   - Handle missing or invalid arguments gracefully

2. **Security**
   - Always validate input
   - Don't trust external data
   - Handle errors appropriately

3. **Performance**
   - Minimize memory usage
   - Process and release resources efficiently
   - Consider timeouts

4. **Debugging**
   - Use informative logs in stderr
   - Include timestamps in logs
   - Structure error messages

5. **Best Practices**
   - Document expected input/output format
   - Keep function stateless
   - Follow idempotency principles