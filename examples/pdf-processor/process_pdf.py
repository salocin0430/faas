import os
import json
import sys
import PyPDF2
import requests
import io

def main(input_data):
    try:
        # Si hay input, es el primer argumento
        input_data = sys.argv[1] if len(sys.argv) > 1 else ""
        
        # Parsear input
        input_json = json.loads(input_data)
        
        # Obtener referencia al PDF
        pdf_ref = input_json["object_inputs"]["pdf_file"]
        pages = input_json["direct_inputs"]["pages"]
        
        # Obtener el archivo usando la URL base desde env
        api_base_url = os.getenv('API_BASE_URL')
        url = f"{api_base_url}/{pdf_ref}"
        
        response = requests.get(url)
        if response.status_code != 200:
            raise Exception(f"Error obteniendo archivo: {response.status_code}")
        
        # Procesar PDF desde los bytes recibidos
        pdf_bytes = io.BytesIO(response.content)
        pdf_reader = PyPDF2.PdfReader(pdf_bytes)
        num_pages = len(pdf_reader.pages)
        text = ""
        
        # Extraer texto de las páginas especificadas
        for page_num in range(min(pages, num_pages)):
            text += pdf_reader.pages[page_num].extract_text()
        
        return json.dumps({
            "pages_processed": min(pages, num_pages),
            "total_pages": num_pages,
            "text": text[:1000]  # Primeros 1000 caracteres
        })
        
    except Exception as e:
        return json.dumps({"error": str(e)}) 