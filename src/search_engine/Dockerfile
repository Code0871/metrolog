## пока что пусто
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bin/server ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/server .
COPY --from=builder /app/config ./config

RUN ls -la /app/config/ && cat /app/config/config.env

EXPOSE 8090

CMD ["/app/server"]
