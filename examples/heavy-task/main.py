import sys
import json
import time
import random

def heavy_processing(seconds):
    # Simular trabajo pesado
    time.sleep(seconds)
    return f"Processed for {seconds} seconds"

def main():
    # Leer input
    input_data = sys.argv[1] if len(sys.argv) > 1 else "{}"
    try:
        data = json.loads(input_data)
        # Tiempo de procesamiento entre 10 y 30 segundos
        process_time = data.get("process_time", random.randint(10, 30))
    except json.JSONDecodeError:
        process_time = 15

    # Ejecutar procesamiento
    result = heavy_processing(process_time)
    
    # Retornar resultado
    print(json.dumps({
        "result": result,
        "processed_for": process_time
    }))

if __name__ == "__main__":
    main() 