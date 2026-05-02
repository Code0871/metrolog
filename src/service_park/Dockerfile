# ============================================================================
# Dockerfile для Go сервиса
# ============================================================================

# Сборка базового образа с Go
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем сервис
RUN go build -o service_park .


# Финальный образ 
FROM alpine:latest

WORKDIR /app

# Копируем собранный бинарный файл из этапа сборки
COPY --from=builder /app/service_park .

# Открытие порта для сервиса
EXPOSE 8080

# Команда для запуска сервиса
CMD ["./service_park"]