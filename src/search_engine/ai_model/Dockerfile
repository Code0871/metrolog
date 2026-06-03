FROM python:3.12-slim

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /app

# RUN pip config set global.index-url https://mirror.yandex.ru/pypi/simple/
# RUN pip config set global.trusted-host mirror.yandex.ru

RUN pip config set global.timeout 100

# зависимости для библиотек
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

COPY pyproject.toml .

RUN pip install --no-cache-dir --upgrade pip     
RUN pip install --no-cache-dir .

# Скачиваем модели
COPY src/download_models.py .
RUN python download_models.py

# Копируем исходный код
COPY src/ ./src/
COPY main.py .

# Создаем директорию для конфига (будет подмонтирована)
RUN mkdir -p /app/config

EXPOSE 8000

CMD ["python", "main.py"]