# Service-specific Dockerfile (standalone multi-stage)
FROM golang:1.23.3-alpine AS builder

ARG SERVICE_NAME=api-gateway

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go.mod for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build service binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /bin/api-gateway \
    ./services/api-gateway/

FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /bin/api-gateway ./bin/api-gateway

# Healthcheck for API Gateway
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:8080/api/v1/health || exit 1

EXPOSE 8080

CMD ["./bin/api-gateway"]