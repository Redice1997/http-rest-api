FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Установите необходимые пакеты для сборки (если нужны для Redis/Prometheus клиентов)
RUN apk add --no-cache gcc musl-dev

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o main ./cmd/api/main.go

FROM alpine:3.18

# Установите CA certificates и временную зону
RUN apk add --no-cache ca-certificates tzdata

# Создайте non-root пользователя для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Скопируйте бинарник из builder stage
COPY --from=builder --chown=appuser:appgroup /build/main /app/

# Переключитесь на non-root пользователя
USER appuser

EXPOSE 8080

# Добавьте healthcheck для мониторинга
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/main"]
