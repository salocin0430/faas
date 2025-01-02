FROM python:3.9-slim

WORKDIR /app
COPY process_pdf.py .

RUN pip install PyPDF2 requests

ENV API_BASE_URL="http://api:9080/api"

ENTRYPOINT ["python", "-c", "import sys, json; from process_pdf import main; print(main(sys.argv[1]))"] 