# Multi-stage Dockerfile base for Go microservices
# Usage in service Dockerfiles: 
# ARG SERVICE_NAME=card-generator
# Then build with docker build --build-arg SERVICE_NAME=card-generator -t card-generator .

ARG SERVICE_NAME=cardgen

FROM golang:1.23.3-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /bin/${SERVICE_NAME} \
    ./cmd/${SERVICE_NAME}/

FROM alpine:3.20 AS runtime

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

# Copy binary
COPY --from=builder /bin/${SERVICE_NAME} ./bin/${SERVICE_NAME}

# Default healthcheck (services override with correct port)
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:${API_GATEWAY_PORT:-8080}/health || exit 1

# Expose configurable port
EXPOSE ${API_GATEWAY_PORT:-8080}

# Run service
CMD ./bin/${SERVICE_NAME}