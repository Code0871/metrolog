# Dockerfile
FROM python:3.11-slim

WORKDIR /app

RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

COPY pyproject.toml .

# Устанавливаем Python-зависимости
RUN pip install --no-cache-dir \
    flask \
    sentence-transformers \
    scikit-learn \
    numpy \
    python-dotenv \
    huggingface-hub

# Копируем исходный код
COPY src/ ./src/
COPY main.py .
COPY model_service.py .


# Копируем модели (если они уже скачаны)
COPY models_cache/ ./models_cache/

# Копируем BM25 индекс
COPY bm25_index.pkl .

# Создаем директорию для конфига (будет подмонтирована)
RUN mkdir -p /app/config

# Переменные окружения
ENV PYTHONUNBUFFERED=1
ENV MODEL_CACHE=/app/models_cache

EXPOSE 8000

CMD ["python", "main.py"]