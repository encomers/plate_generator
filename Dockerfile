# syntax=docker/dockerfile:1.4

# =============================================================================
# Stage 1: Builder
# =============================================================================
FROM golang:1.25-alpine3.21 AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git ca-certificates make bash

# Устанавливаем рабочую директорию
WORKDIR /build

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости (кэшируется при неизменных go.mod/go.sum)
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Копируем исходный код проекта
COPY . .

# Генерируем Swagger-документацию
# Аннотации находятся в server.go, main.go в cmd/generator/
RUN --mount=type=cache,target=/go/pkg/mod \
    go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init -g cmd/generator/main.go -o ./docs 2>/dev/null || echo "Swagger generation skipped"

# Собираем бинарный файл с оптимизациями
# CGO_ENABLED=0 для статической линковки, -trimpath для воспроизводимых сборок
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-w -s -extldflags '-static'" \
    -o /bin/generator ./cmd/generator

# =============================================================================
# Stage 2: Runtime (минимальный образ)
# =============================================================================
FROM alpine:3.19 AS runtime

# Устанавливаем корневые сертификаты для HTTPS-запросов и tzdata для времени
RUN apk --no-cache add ca-certificates tzdata

# Создаём не-привилегированного пользователя для безопасности
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Рабочая директория
WORKDIR /app

# Копируем бинарный файл из stage builder
COPY --from=builder --chown=appuser:appgroup /bin/generator /app/generator

# Переключаемся на не-привилегированного пользователя
USER appuser

# Экспортируем порт (по умолчанию 8080)
EXPOSE 8080

# Переменные окружения по умолчанию
ENV GIN_MODE=release \
    PORT=8080 \
    LOG_LEVEL=info

# Точка входа
ENTRYPOINT ["/app/generator"]
