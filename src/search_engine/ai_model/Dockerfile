FROM python:3.10-slim

WORKDIR /app

RUN apt-get update && apt-get install -y \
    gcc g++ \
    && rm -rf /var/lib/apt/lists/*

RUN pip install --upgrade pip

COPY pyproject.toml .

RUN pip install --no-cache-dir --root-user-action=ignore .

COPY model_service.py .

EXPOSE 8000

CMD ["python", "model_service.py"]