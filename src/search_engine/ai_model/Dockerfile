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
    huggingface-hub \
    rank_bm25==0.2.2

# Копируем исходный код
COPY src/ ./src/
COPY main.py .
COPY model_service.py .

# Копируем модели (если они уже скачаны)
COPY models_cache/ ./models_cache/

# Создаем директорию для конфига (будет подмонтирована)
RUN mkdir -p /app/config

# Переменные окружения
ENV PYTHONUNBUFFERED=1
ENV MODEL_CACHE=/app/models_cache

EXPOSE 8000

CMD ["python", "main.py"]