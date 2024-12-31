import sys
import json

def main():
    # Si hay input, es el primer argumento
    input_data = sys.argv[1] if len(sys.argv) > 1 else ""
    
    try:
        # Intentar parsear como JSON
        data = json.loads(input_data) if input_data else {}
        message = data.get("message", "Hello World!")
    except json.JSONDecodeError:
        # Si no es JSON válido, usar el input como string
        message = input_data or "Hello World!"
    
    print(message)

if __name__ == "__main__":
    main() 