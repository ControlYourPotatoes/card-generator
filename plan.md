# Card Generator Modernization Plan

## Docker Containerization & Microservices Architecture

### IMPORTANT AGENT RULES

🚨 **DO NOT PROCEED TO THE NEXT PHASE WITHOUT EXPLICIT USER APPROVAL**

- Complete current phase fully before requesting permission to continue
- Ask user to review and approve each phase completion
- Wait for user confirmation before starting next phase
- Test thoroughly between phases

---

## Project Overview

**Objective**: Modernize the card-generator project with Docker containerization and functional microservices architecture while maintaining the solid domain-driven design already in place. Enable full-stack demonstration at new employer showing both product capability (card generation pipeline) and technical sophistication (containerized microservices).

**Current Architecture Strengths**:

- Clean separation of concerns with well-defined interfaces (Store, CardGenerator, ArtProcessor, TextProcessor)
- Domain-driven design with core models isolated from implementation details
- Dependency injection container enabling loose coupling
- Comprehensive test coverage and build automation
- Full-stack setup ready for enhancement (Next.js frontend + Go backend)

**Modernization Goals**:

1. Containerize entire application stack with Docker
2. Extract backend into functional microservices (Card Generator, Image Renderer, Database, API Gateway)
3. Maintain API compatibility with existing frontend
4. Enable local development with docker-compose
5. Create production-ready deployment configuration
6. Prepare for future scaling and feature additions
7. Integrate Open Router connectivity for agentic development workflows

---

## Phase 1: Docker Foundation & Service Architecture

### Phase 1.1: Dockerfile & Docker Configuration Strategy

**Objectives**:

- Create optimized multi-stage Dockerfiles for backend services
- Set up docker-compose for local development environment
- Establish service communication patterns
- Create Docker networking configuration

**Tasks**:

#### 1.1.1 Backend Base Dockerfile

Location: `/backend/Dockerfile`

Create multi-stage build:

- Stage 1: Build stage (Go 1.23.3, compile binaries)
- Stage 2: Runtime stage (minimal alpine image)
- Build binaries for: cardgen-di, importer, db-cleaner as service entrypoints

**Key considerations**:

- Build flags for optimization: `-ldflags "-s -w"`
- Include health checks
- Set proper working directories
- Base image: `golang:1.23.3-alpine` for build, `alpine:latest` for runtime

#### 1.1.2 Service-Specific Docker Configurations

Location: `/backend/services/` (new directory structure)

Create service directories that will become individual microservices:

- `/backend/services/card-generator/Dockerfile` - Card generation service
- `/backend/services/image-renderer/Dockerfile` - Image rendering service
- `/backend/services/api-gateway/Dockerfile` - API routing service
- `/backend/services/importer/Dockerfile` - CSV import service

Each service Dockerfile:

- Uses shared backend base layer (multi-stage inheritance pattern)
- Has single entrypoint/responsibility
- Includes service-specific configuration

#### 1.1.3 docker-compose.yml

Location: `/docker-compose.yml` (root)

Services to define:

- `postgres` - PostgreSQL 16 database, volumes for data persistence
- `card-generator` - Card generation microservice, depends_on postgres
- `image-renderer` - Image rendering microservice, depends_on postgres
- `api-gateway` - API routing service, ports 8080:8080, depends_on others
- `frontend` - Next.js application, ports 3000:3000, depends_on api-gateway
- `adminer` (optional) - Database admin tool, ports 8081:8080

**Network**: Create named network `card-gen-network` with custom bridge driver

#### 1.1.4 Database Service Configuration

Location: `/docker-compose.yml` postgres service section

PostgreSQL configuration:

- Image: `postgres:16-alpine`
- Environment: POSTGRES_DB, POSTGRES_USER, POSTGRES_PASSWORD (from .env.docker)
- Volume: `postgres_data:/var/lib/postgresql/data`
- Health check: pg_isready command
- Initialization: Scripts in `/backend/scripts/db-init/` directory

#### 1.1.5 Environment Configuration for Docker

Location: `.env.docker` (root)

Define environment variables for Docker services:

```text
DATABASE_URL=postgresql://cardgen_user:password@postgres:5432/card_generator
POSTGRES_USER=cardgen_user
POSTGRES_PASSWORD=password
POSTGRES_DB=card_generator
API_GATEWAY_PORT=8080
FRONTEND_PORT=3000
GO_ENV=docker
```

**Phase 1.1 Acceptance Criteria**:

- [ ] Multi-stage Dockerfile compiles Go binaries without errors
- [ ] Docker images build successfully without warnings
- [ ] docker-compose.yml is syntactically valid and passes `docker-compose config`
- [ ] PostgreSQL container starts and health check passes
- [ ] All services defined with proper dependencies and networking
- [ ] Named network `card-gen-network` enables service-to-service communication
- [ ] Volume configuration preserves database data across container restarts

**Testing Required Before Phase 1.2**:

- Build all Docker images: `docker-compose build`
- Start services: `docker-compose up -d`
- Verify PostgreSQL is accessible: `docker-compose exec postgres psql -U cardgen_user -d card_generator -c "SELECT 1"`
- Check inter-service network connectivity
- Confirm services are running: `docker-compose ps`
- Stop and verify data persistence: `docker-compose down` then `docker-compose up -d` and check database

### 🛑 STOP: Request user approval before proceeding to Phase 1.2

---

### Phase 1.2: API Gateway & REST Service Contracts

**Objectives**:

- Define REST API contracts for frontend ↔ API Gateway communication
- Create API Gateway service to route HTTP requests from frontend
- Establish REST endpoints for image serving and data retrieval
- Maintain tRPC type safety on frontend

**Architecture Pattern - Hybrid Approach**:

```text
Frontend (Next.js) ← → tRPC (type-safe) ← → API Gateway (REST/HTTP)
                                              ├─ card-generator:8081
                                              ├─ image-renderer:8082  (HTTP for image streaming)
                                              └─ importer:8084
                                    ↓
                     gRPC (internal service-to-service, Phase 1.5)
                                    ↓
                    card-generator ← → image-renderer
                         ↓ ↓ ↓ (gRPC contracts)
                        (type-safe service boundaries)
```

**Tasks**:

#### 1.2.1 Service Contract Definitions

Location: `/api/contracts.ts` (update existing types)

Define REST API contracts for frontend consumption:

- `POST /api/v1/cards/generate` - Generate card from JSON data
- `GET /api/v1/cards/:id` - Retrieve card metadata
- `GET /api/v1/cards/:id/render` - Stream card image (PNG/SVG)
- `POST /api/v1/cards/analyze` - Analyze card for tags and metadata
- `POST /api/v1/import/csv` - Import cards from CSV (async job)
- `GET /api/v1/import/:jobId/status` - Check import job status
- `DELETE /api/v1/admin/cards` - Clear all cards (admin endpoint)
- `GET /api/v1/health` - Service health status

Include request/response types with Zod validation. Version API as `/api/v1/` for future compatibility.

#### 1.2.2 API Gateway Service Implementation

Location: `/backend/services/api-gateway/main.go` (new)

Create HTTP routing layer:

- Use `github.com/go-chi/chi` for routing (better than net/http for microservices)
- Routes to backend services via docker network DNS (`card-generator:8081`, `image-renderer:8082`, etc.)
- Request/response logging middleware with correlation IDs
- Image streaming handler for binary responses (PNG/SVG)
- Health check aggregation endpoint: `GET /api/v1/health`
- CORS headers for frontend communication
- API versioning support (`/api/v1/*`)

Key routes:

```go
POST   /api/v1/cards/generate      → card-generator:8081/generate (JSON)
GET    /api/v1/cards/:id           → card-generator:8081/card/:id (JSON)
GET    /api/v1/cards/:id/render    → image-renderer:8082/render/:id (image stream)
POST   /api/v1/cards/analyze       → card-generator:8081/analyze (JSON)
POST   /api/v1/import/csv          → importer:8084/import (multipart/form-data)
GET    /api/v1/import/:jobId       → importer:8084/status/:jobId (JSON)
DELETE /api/v1/admin/cards         → card-generator:8081/admin/clear (JSON)
GET    /api/v1/health              → Aggregate all service health (JSON)
```

#### 1.2.3 Update Frontend tRPC Router

Location: `/frontend/src/server/api/routers/` (modify existing)

Update tRPC procedures to call `http://api-gateway:8080/api/v1/*`:

- Proper error handling for service communication failures (timeouts, connection errors)
- Timeout configuration: 30s for card operations, 60s for imports, 5s for health checks
- Request/response type compatibility with Zod schemas
- Automatic retry logic for transient failures (429, 503)
- Request correlation IDs for debugging across services

**Phase 1.2 Acceptance Criteria**:

- [ ] API contracts defined and validated with Zod types
- [ ] REST API uses versioning (`/api/v1/*`)
- [ ] API Gateway HTTP router defined with all service routes
- [ ] All routes return proper HTTP status codes (200, 201, 400, 404, 500)
- [ ] Image streaming works correctly for binary PNG/SVG responses
- [ ] Content-Type headers correct for different response types
- [ ] CORS headers properly configured for frontend access
- [ ] Health check aggregates all service status
- [ ] Logger middleware captures requests with correlation IDs
- [ ] tRPC router updated to use API Gateway endpoint with proper error handling
- [ ] Retry logic implements exponential backoff

**Testing Required Before Phase 1.3**:

- Start docker-compose stack: `docker-compose up -d`
- Call API Gateway health: `curl http://localhost:8080/api/v1/health`
- Test service discovery by names from within containers
- Verify frontend can communicate with api-gateway
- Test image streaming: `curl http://localhost:8080/api/v1/cards/<id>/render -o test.png`
- Verify binary image data is valid: `file test.png`
- Test JSON endpoints with curl/Postman
- Check logs for correlation IDs: `docker-compose logs api-gateway | grep "request-id"`
- Test error scenarios: service timeouts, 404s, 500s

### 🛑 STOP: Request user approval before proceeding to Phase 1.3

---

### Phase 1.3: Observability & Error Handling Foundation

**Objectives**:

- Establish structured logging across all services
- Create standardized error types with context tracking
- Implement request correlation IDs for debugging
- Set up health check patterns for all services

**Tasks**:

#### 1.3.1 Structured Logging System

Location: `/backend/pkg/logging/logger.go` (new)

Implement structured logging interface:

```go
type Logger interface {
  Debug(ctx context.Context, msg string, fields ...interface{})
  Info(ctx context.Context, msg string, fields ...interface{})
  Warn(ctx context.Context, msg string, fields ...interface{})
  Error(ctx context.Context, msg string, err error, fields ...interface{})
}
```

Features:

- JSON output for log aggregation (parseable by docker logs)
- Automatic correlation ID inclusion in all logs
- Request/response logging middleware
- Performance timing for operations
- Log levels configurable per environment (.env)

#### 1.3.2 Error Handling with Context

Location: `/backend/pkg/errors/errors.go` (new)

Define standardized error types:

```go
type Error struct {
  Code      string                 // Machine-readable error code (e.g., "CARD_NOT_FOUND")
  Message   string                 // User-friendly message
  Details   map[string]interface{} // Additional context
  Cause     error                  // Underlying error (if any)
  HTTPCode  int                    // HTTP status code
  RequestID string                 // For correlation
}
```

Standard error codes:

- `INVALID_INPUT` - Validation errors (400)
- `NOT_FOUND` - Resource not found (404)
- `CONFLICT` - Business logic conflict (409)
- `SERVICE_ERROR` - Internal service error (500)
- `SERVICE_UNAVAILABLE` - Dependent service down (503)
- `TIMEOUT` - Service timeout (504)

#### 1.3.3 Request Correlation IDs

Location: `/backend/pkg/middleware/correlation.go` (new)

Middleware implementation:

- Generate UUID for each request if not provided
- Include in all request logs
- Pass to downstream services in headers: `X-Request-ID`
- Return in response headers for client reference

#### 1.3.4 Health Check Pattern

Location: `/backend/pkg/health/checks.go` (new)

Standard health check interface:

```go
type HealthCheck interface {
  Check(ctx context.Context) *HealthStatus
}

type HealthStatus struct {
  Status      string                 // "ok", "degraded", "error"
  Checks      map[string]CheckStatus // Individual component status
  Timestamp   time.Time
  Version     string
  Uptime      time.Duration
}
```

Each service implements:

- Database connectivity check
- Dependent service connectivity check
- Memory/resource usage check
- Graceful degradation when dependencies unavailable

**Phase 1.3 Acceptance Criteria**:

- [ ] Structured logging implemented and working in all services
- [ ] All log output is valid JSON
- [ ] Correlation IDs visible in logs
- [ ] Standardized error types defined
- [ ] Error responses include proper HTTP codes
- [ ] Health check endpoint works for all services
- [ ] Health check includes dependency status
- [ ] Request timing appears in logs
- [ ] Error context captured and loggable without data loss

**Testing Required Before Phase 1.4**:

- Start one service: `docker-compose up card-generator`
- Generate error (e.g., invalid input)
- Check logs are valid JSON: `docker-compose logs card-generator | jq .`
- Verify correlation ID in logs
- Test health endpoint: `curl http://localhost:8081/health | jq .`
- Kill a dependency (simulate DB down) and check health status
- Verify error responses include codes and details

### 🛑 STOP: Request user approval before proceeding to Phase 1.4

---

### Phase 1.4: Event Bus Foundation for Async Operations

**Objectives**:

- Establish event-driven architecture for image processing
- Decouple card generation from image rendering
- Enable async operations and progress tracking
- Create event system for future domain events

**Tasks**:

#### 1.4.1 Event Bus Interface

Location: `/backend/pkg/events/bus.go` (new)

Define event system:

```go
type Event interface {
  EventType() string
  EventID() string
  Timestamp() time.Time
  AggregateID() string // Correlate related events
}

type EventBus interface {
  Publish(ctx context.Context, event Event) error
  Subscribe(eventType string, handler EventHandler) error
}

type EventHandler interface {
  Handle(ctx context.Context, event Event) error
}
```

#### 1.4.2 Domain Events

Location: `/backend/internal/events/card_events.go` (new)

Define card-related events:

```go
type CardGenerationRequested struct {
  EventID    string
  CardID     string
  CardData   *card.CardDTO
  Format     string // "png", "svg", or "both"
  Priority   int
  Timestamp  time.Time
}

type CardImage GenerationCompleted struct {
  EventID    string
  CardID     string
  ImageURL   string
  Format     string
  Status     string // "success", "failed"
  Error      string
  Timestamp  time.Time
}

type CardImageGenerationFailed struct {
  EventID    string
  CardID     string
  Format     string
  Error      string
  Timestamp  time.Time
}
```

#### 1.4.3 In-Process Event Bus Implementation

Location: `/backend/pkg/events/inmemory/bus.go` (new)

Development implementation:

- In-memory queue for local testing
- Synchronous for simplicity
- No persistence (for demo purposes)
- Can be swapped for message queue (RabbitMQ, Redis) later

#### 1.4.4 Image Generation Event Listener

Location: `/backend/services/image-renderer/listener.go` (new)

Subscribe to `CardImageGenerationRequested`:

- Receive generation request
- Render image asynchronously
- Publish `CardImageGenerationCompleted` on success
- Publish `CardImageGenerationFailed` on error
- Update progress in job tracking system

#### 1.4.5 Integration with Card Generator Service

Location: `/backend/services/card-generator/handlers/generate.go` (modify)

Update card generation handler:

- After card stored, publish `CardImageGenerationRequested` event
- Return card immediately (don't wait for image)
- Include image URL in response once image completes (via event)
- Handle event callback to update card metadata

**Phase 1.4 Acceptance Criteria**:

- [ ] Event bus interface defined and documented
- [ ] Domain events defined for card operations
- [ ] In-memory event bus implementation working
- [ ] Image generation decoupled from card creation
- [ ] Events published and consumed correctly
- [ ] Event subscribers handle success and failure cases
- [ ] Correlation IDs flow through events
- [ ] Event handler errors logged properly

**Testing Required Before Phase 1.5**:

- Start card-generator and image-renderer services
- Generate card via API
- Verify `CardImageGenerationRequested` event published
- Check image-renderer receives and processes event
- Verify `CardImageGenerationCompleted` event published
- Check correlation IDs propagate through events
- Simulate image renderer failure and verify error event
- Monitor logs for complete event flow

### 🛑 STOP: Request user approval before proceeding to Phase 1.5

---

### Phase 1.5: gRPC Service Contracts for Internal Communication

**Objectives**:

- Define gRPC contracts between microservices (Go type-safe)
- Replace internal HTTP calls with gRPC for performance
- Establish service registry for discovery
- Create code generation from Protocol Buffers

**Tasks**:

#### 1.5.1 Protocol Buffer Definitions

Location: `/backend/proto/card_service.proto` (new)

Define service contracts:

```protobuf
syntax = "proto3";
package cardgenerator;

service CardGenerator {
  rpc GenerateCard(GenerateCardRequest) returns (GenerateCardResponse);
  rpc GetCard(GetCardRequest) returns (CardData);
  rpc ListCards(ListCardsRequest) returns (ListCardsResponse);
  rpc AnalyzeCard(AnalyzeCardRequest) returns (AnalysisResult);
}

service ImageRenderer {
  rpc RenderCardImage(RenderRequest) returns (RenderResponse);
  rpc GetImageStatus(ImageStatusRequest) returns (ImageStatusResponse);
}

message CardData {
  string id = 1;
  string name = 2;
  int32 cost = 3;
  string card_type = 4;
  string effect = 5;
  repeated string keywords = 6;
  map<string, string> metadata = 7;
}

message GenerateCardRequest {
  CardData card = 1;
  string request_id = 2; // Correlation ID
}

message GenerateCardResponse {
  CardData card = 1;
  string status = 2;
  string error = 3;
}

// ... additional messages
```

#### 1.5.2 Code Generation Setup

Location: `/backend/Makefile` (update)

Add protobuf compilation:

```makefile
.PHONY: proto-generate
proto-generate:
  protoc --go_out=. --go-grpc_out=. ./proto/*.proto
```

Dependencies:

- `google.golang.org/protobuf` - Protocol Buffer library
- `google.golang.org/grpc` - gRPC framework
- `protoc` compiler

#### 1.5.3 gRPC Service Implementation

Location: `/backend/services/card-generator/grpc/server.go` (new)

Implement gRPC server:

- Separate port: 5051 (internal only, not exposed to frontend)
- Implement generated interfaces from Protocol Buffers
- Reuse existing business logic (generators, stores, analyzers)
- Proper error handling with gRPC error codes
- Context propagation for correlation IDs

Example:

```go
type CardGeneratorServer struct {
  service *CardService
}

func (s *CardGeneratorServer) GenerateCard(ctx context.Context, req *pb.GenerateCardRequest) (*pb.GenerateCardResponse, error) {
  // Extract correlation ID from context
  correlationID := middleware.CorrelationIDFromContext(ctx)

  // Call existing service
  result, err := s.service.GenerateCard(ctx, req.Card)
  // ...
}
```

#### 1.5.4 gRPC Client Implementation

Location: `/backend/services/image-renderer/grpc/client.go` (new)

Create gRPC clients for service discovery:

- Connect to card-generator:5051 via DNS
- Connection pooling and keepalive
- Automatic retry with backoff
- Timeout configuration per method

#### 1.5.5 Service Registry Pattern

Location: `/backend/pkg/registry/service.go` (new)

Service discovery:

- Services register at startup with startup metadata
- Discovery via docker DNS (no external registry needed for docker-compose)
- Health checks on connections
- Automatic reconnection on failure

**Phase 1.5 Acceptance Criteria**:

- [ ] Protocol Buffers defined for all service contracts
- [ ] Code generation produces Go code without errors
- [ ] gRPC services implemented and compilable
- [ ] gRPC clients created for service-to-service calls
- [ ] Service discovery works via docker DNS
- [ ] Correlation IDs propagate through gRPC calls
- [ ] Error codes mapped from gRPC to standard errors
- [ ] Connection pooling working
- [ ] Timeouts configured per service method
- [ ] API Gateway still routes REST for frontend

**Testing Required Before Phase 2**:

- Compile protobuf: `make proto-generate`
- Start all services: `docker-compose up -d`
- Call card-generator via gRPC (use grpcurl for testing)
- Generate card, verify image-renderer receives via gRPC
- Check performance: gRPC vs previous HTTP calls
- Verify correlation IDs in logs across gRPC boundaries
- Test gRPC connection pooling
- Simulate service failure and verify reconnection
- Monitor gRPC traffic: `docker-compose exec card-generator ss -tulnp | grep 5051`

### 🛑 STOP: Request user approval before proceeding to Phase 2

---

---

## Phase 2: Microservices Decomposition

### Phase 2.1: Card Generator Service Extraction

**Objectives**:

- Extract card generation logic into independent microservice
- Create service endpoints for card CRUD operations
- Implement database persistence for generated cards

**Tasks**:

#### 2.1.1 Card Generator Service Structure

Location: `/backend/services/card-generator/` (new directory)

Structure:

```text
/backend/services/card-generator/
├── main.go                    # Service entry point, HTTP handlers
├── config/
│   └── config.go             # Service-specific config
├── handlers/
│   ├── generate.go           # POST /generate - card generation
│   ├── get.go                # GET /card/:id - retrieve card
│   ├── list.go               # GET /cards - list all cards
│   └── health.go             # GET /health - service health
└── service.go                # Business logic orchestration
```

#### 2.1.2 Service Main Implementation

Location: `/backend/services/card-generator/main.go`

Implement HTTP server:

- Port: 8081
- Health check endpoint
- Graceful shutdown handling
- Dependency initialization (Store, Generator, Parser, Tagger from bootstrap)

Key handlers:

```go
POST /generate
- Request: CardData JSON
- Response: { id, name, type, tags, metadata }
- Process: Parse → Store → Generate image → Analyze → Return

GET /card/:id
- Response: Card metadata with image path

GET /cards
- Response: List of all cards with pagination

GET /health
- Response: { status: "ok" }
```

#### 2.1.3 Database Integration

Location: `/backend/services/card-generator/service.go`

Configure PostgreSQL store:

- Use DATABASE_URL environment variable
- Initialize schema on startup
- Connection pooling configuration
- Migration handling

#### 2.1.4 Update CLI Tool Reference

Location: `/backend/cmd/cardgen-di/main.go`

Keep in place for backward compatibility but:

- Add deprecation note in README
- Redirect users to containerized service
- Plan removal in future version

**Phase 2.1 Acceptance Criteria**:

- [ ] Card Generator service compiles without errors
- [ ] Service starts on port 8081 and responds to health check
- [ ] All handler endpoints return correct HTTP status codes
- [ ] Cards are persisted to PostgreSQL database
- [ ] Service can be called from API Gateway
- [ ] Image generation works within containerized environment
- [ ] Error responses include proper error messages

**Testing Required Before Phase 2.2**:

- Start service: `docker-compose up card-generator -d`
- Health check: `curl http://localhost:8081/health`
- Generate card: `curl -X POST http://localhost:8081/generate -d '{"name":"Test","cost":3,...}'`
- Retrieve card: `curl http://localhost:8081/card/{id}`
- Verify database rows: `docker-compose exec postgres psql -c "SELECT * FROM cards"`
- Check service logs: `docker-compose logs card-generator`

### 🛑 STOP: Request user approval before proceeding to Phase 2.2

---

### Phase 2.2: Image Renderer Service Extraction

**Objectives**:

- Extract image rendering (PNG/SVG) into dedicated microservice
- Implement high-performance rendering with caching
- Support multiple output formats

**Tasks**:

#### 2.2.1 Image Renderer Service Structure

Location: `/backend/services/image-renderer/` (new directory)

Structure:

```text
/backend/services/image-renderer/
├── main.go                    # Service entry point
├── handlers/
│   ├── render.go             # POST /render/:cardId - render card image
│   └── health.go             # GET /health
├── cache/
│   └── image_cache.go        # Local image caching
└── config/
    └── config.go             # Service config
```

#### 2.2.2 Renderer Service Implementation

Location: `/backend/services/image-renderer/main.go`

HTTP server:

- Port: 8082
- Handlers for rendering endpoints
- Cache management
- Graceful shutdown

Key handler:

```go
POST /render/:cardId
- Path param: cardId (string)
- Query param: format (png|svg, default: png)
- Response: Binary image data with Content-Type
- Process: Fetch card data → Render → Cache → Serve
```

#### 2.2.3 Image Caching Implementation

Location: `/backend/services/image-renderer/cache/image_cache.go`

File-based caching:

- Cache directory: `/tmp/card-images/`
- Filename: `{cardId}.{format}`
- TTL: 24 hours (configurable)
- Cache invalidation on card updates

#### 2.2.4 Integration with Card Generator

Location: `/backend/services/card-generator/handlers/generate.go`

Update generate handler:

- After card is stored, call image-renderer service
- Include rendered image URL in response
- Async rendering option for large batches

**Phase 2.2 Acceptance Criteria**:

- [ ] Image Renderer service compiles and starts on port 8082
- [ ] Render endpoint returns valid PNG/SVG image data
- [ ] Content-Type headers are correct for image format
- [ ] Image caching works and improves performance on cache hits
- [ ] Service handles missing cards gracefully (returns 404)
- [ ] Cache invalidation works correctly

**Testing Required Before Phase 2.3**:

- Start service: `docker-compose up image-renderer -d`
- Generate test card via card-generator service
- Call render endpoint: `curl http://localhost:8082/render/{cardId}?format=png -o test.png`
- Verify PNG file is valid: `file test.png`
- Check cache directory: `docker-compose exec image-renderer ls -la /tmp/card-images/`
- Call render again, verify cache hit improves response time
- Test SVG format: `curl http://localhost:8082/render/{cardId}?format=svg`

### 🛑 STOP: Request user approval before proceeding to Phase 2.3

---

### Phase 2.3: CSV Importer Service Extraction

**Objectives**:

- Extract CSV import logic into dedicated microservice
- Enable bulk card import without CLI
- Support async import with progress tracking

**Tasks**:

#### 2.3.1 Importer Service Structure

Location: `/backend/services/importer/` (new directory)

Structure:

```text
/backend/services/importer/
├── main.go                    # Service entry point
├── handlers/
│   ├── import.go             # POST /import - CSV import handler
│   ├── status.go             # GET /import/:jobId - import status
│   └── health.go             # GET /health
├── worker/
│   └── import_worker.go       # Background import jobs
└── config/
    └── config.go
```

#### 2.3.2 Importer Service Implementation

Location: `/backend/services/importer/main.go`

HTTP server:

- Port: 8084
- Handlers for import operations
- Job tracking support
- Graceful shutdown with job cleanup

Key handlers:

```go
POST /import
- Body: multipart/form-data with CSV file
- Query: cardType (creature, spell, etc.), dryRun (bool)
- Response: { jobId, status: "started" }
- Process: Queue job → Return jobId immediately

GET /import/:jobId
- Response: { jobId, status, importedCount, totalCount, errors[] }
- Process: Return current job status
```

#### 2.3.3 Background Job Processing

Location: `/backend/services/importer/worker/import_worker.go`

Implement worker pool:

- Process import jobs asynchronously
- Track import progress
- Store results and errors
- Persist imported cards to database

#### 2.3.4 Keep CLI Tool for Backward Compatibility

Location: `/backend/cmd/importer/main.go`

- Maintain existing functionality
- Document that service endpoint is recommended route
- Plan removal in future version

**Phase 2.3 Acceptance Criteria**:

- [ ] Importer service compiles and starts on port 8084
- [ ] CSV upload endpoint accepts multipart file uploads
- [ ] Multiple card types supported (creature, spell, artifact, etc.)
- [ ] Dry-run mode validates without persisting
- [ ] Import job tracking returns accurate status
- [ ] Imported cards appear in database
- [ ] Error handling for invalid CSV format

**Testing Required Before Phase 3**:

- Start service: `docker-compose up importer -d`
- Upload CSV: `curl -X POST -F "file=@test.csv" -F "cardType=creature" http://localhost:8084/import`
- Check job status: `curl http://localhost:8084/import/{jobId}`
- Verify cards in database: `docker-compose exec postgres psql -c "SELECT COUNT(*) FROM cards"`
- Test with invalid CSV, verify error handling

### 🛑 STOP: Request user approval before proceeding to Phase 3

---

### Phase 3: Full Integration & Demo Readiness

**Objectives**:

- Connect all services through API Gateway
- Verify end-to-end card generation pipeline
- Enable frontend to leverage all microservices
- Prepare for production deployment
- Set up monitoring and debugging

**Tasks**:

#### 3.1 API Gateway Full Integration

Location: `/backend/services/api-gateway/main.go`

Enhance API Gateway:

- Route all card operations through appropriate microservices
- Implement request/response logging
- Add service health monitoring
- Expose `/api/admin/status` endpoint showing all service health

#### 3.2 Frontend Integration Testing

Location: `/frontend/src/server/api/routers/`

Update tRPC routers to use containerized API:

- Ensure API calls go through `http://api-gateway:8080`
- Implement proper error handling for service timeouts
- Add retry logic for transient failures

#### 3.3 End-to-End Demo Flow

Document and verify:

1. User creates card via frontend form
2. Request routed through tRPC → API Gateway
3. Card Generator service parses and stores card
4. Image Renderer service generates PNG/SVG
5. Frontend displays card with image

#### 3.4 Production Configuration

Location: `/docker-compose.prod.yml` (new)

Create production override:

- Remove development tools (Adminer)
- Use proper secrets management (Docker Secrets or external vault)
- Add resource limits (memory, CPU)
- Configure logging to centralized system
- Enable health checks on all services

#### 3.5 Documentation

Location: `/DEPLOYMENT.md` (new)

Document:

- Local development: `docker-compose up`
- Production deployment steps
- Service configuration options
- Troubleshooting guide
- Architecture diagram

**Phase 3 Acceptance Criteria**:

- [ ] API Gateway healthy and routing requests correctly
- [ ] Frontend successfully communicates with all backend services
- [ ] Complete card generation pipeline works end-to-end
- [ ] Frontend displays generated cards with images
- [ ] All services have proper error handling and logging
- [ ] Service health endpoints accessible
- [ ] Production docker-compose configuration ready
- [ ] Deployment documentation complete

**Testing Required Before Phase 4**:

- Full stack: `docker-compose up`
- Frontend loads: `http://localhost:3000`
- Create card through UI, verify:
  - Data stored in database
  - Image generated and displayed
  - Card metadata visible in API responses
- Load testing with bulk import
- Service restart resilience
- Data persistence across stack restart

### 🛑 STOP: Request user approval before proceeding to Phase 4

---

#### Phase 4: Demo Optimization & Open Router Integration

**Objectives**:

- Optimize card generation for demonstration
- Enable better models via Open Router for development
- Prepare comprehensive demo scenario
- Document all capabilities

**Tasks**:

#### 4.1 Demo Scenario Setup

Location: `/demo/` (new directory)

Create demo data and workflows:

- Sample CSV with compelling cards
- Demo script automating full pipeline
- Pre-generated card examples
- Frontend walkthrough documentation

#### 4.2 Performance Optimization

Location: `/backend/services/` (update all)

Optimize for demo:

- Image caching strategy for fast re-renders
- Batch operation support for importing
- Response time targets: <500ms for card generation
- Memory footprint optimization

#### 4.3 Open Router Integration Setup

Location: `/vscode-settings.json` (new)

Document Open Router setup for VSCode:

- API key configuration
- Model selection (Claude 3.5 Sonnet recommended)
- Endpoint configuration
- Per-project settings

This enables agentic coding workflow while developing future enhancements (Phase 3 template system, etc.)

#### 4.4 Deployment Readiness

Location: `/scripts/deploy.sh` (new)

Create deployment automation:

- Build and push Docker images
- Deploy to target environment
- Run smoke tests
- Verify all services

**Phase 4 Acceptance Criteria**:

- [ ] Demo cards load and render in <500ms
- [ ] Demo workflow completes successfully end-to-end
- [ ] All performance targets met
- [ ] Open Router integration documented and working
- [ ] Deployment script tested and ready
- [ ] Documentation complete for new employer demo

**Testing Required Before Phase 4 Completion**:

- Run demo workflow start-to-finish
- Measure performance metrics
- Verify Open Router connectivity from VSCode
- Test deployment script in clean environment
- Run all services in demo mode
- Screenshot card samples for presentation

### 🛑 COMPLETION: Phase 4 is the final demo-ready phase

---

## Future Phase: SVG Template System (Phase 3 of SVG Migration)

**Status**: Planned for post-demo enhancement
**Dependencies**: Complete Phase 4 first
**Scope**: Implement Phase 3 from SVG_TEMPLATE_SYSTEM_PLAN.md with microservices compatibility

---

## Architecture Compliance Checklist

### SOLID Principles

- [ ] Single Responsibility: Each microservice has one clear responsibility (Card Gen, Image Render, Import)
- [ ] Open/Closed: Services extensible via interface boundaries without modification
- [ ] Liskov Substitution: Store implementations (Memory, PostgreSQL) are interchangeable
- [ ] Interface Segregation: Services only expose necessary endpoints
- [ ] Dependency Inversion: Services depend on abstractions (interfaces), not concrete implementations

### Clean Architecture

- [ ] Dependencies point inward: CLI tools/HTTP handlers depend on services, not vice versa
- [ ] Business logic isolated: Core generation/analysis logic in internal packages
- [ ] Framework-independent: Could swap HTTP framework without changing business logic

### Domain-Driven Design

- [ ] Rich domain models: Card types, validation preserved
- [ ] Ubiquitous language: CardType, CardTagger, ArtProcessor terminology consistent
- [ ] Bounded contexts: Card domain separate from API delivery layer

---

## Risk Mitigation

1. **Service Communication Failures**: Implement circuit breakers and retry logic with exponential backoff; fallback to in-process if service down (graceful degradation)
2. **Database Connectivity**: Use connection pooling, health checks, automatic reconnection; keep migration scripts in version control for recovery
3. **Image Storage**: Local caching with TTL; volume mounting ensures persistence; fallback generation if cache miss
4. **Frontend Breaking Changes**: Maintain API versioning (/api/v1, /v2); deprecation warnings; semantic versioning
5. **Deployment Complexity**: Comprehensive docker-compose.yml covers both dev and prod; automated smoke tests catch issues early
6. **Performance Regression**: Benchmark before/after containerization; cache strategy prevents N+1 renderings

---

## Success Metrics

- [ ] **Functional**: Card generation pipeline works end-to-end through containerized services
- [ ] **Deployment**: Application starts with single `docker-compose up` command
- [ ] **Performance**: Card generation <500ms, image rendering <300ms, API response <200ms
- [ ] **Compatibility**: Frontend works against containerized backend without changes
- [ ] **Scalability**: Services can run independently and be scaled horizontally
- [ ] **Maintainability**: Clear service boundaries, documented APIs, comprehensive logging
- [ ] **Demo Readiness**: Demo scenario runs flawlessly, showcasing product and technical sophistication

---

## Critical Files to Create/Modify

**Phase 1.1 - Docker Foundation**:

- `/Dockerfile` - Multi-stage backend build
- `/docker-compose.yml` - Service orchestration
- `.env.docker` - Docker environment variables

**Phase 1.2 - API Gateway & REST Contracts**:

- `/api/contracts.ts` - REST API contracts with Zod types (using `/api/v1/*` versioning)
- `/backend/services/api-gateway/main.go` - API Gateway service (chi router, REST endpoints)
- `/backend/services/api-gateway/middleware/logging.go` - Request logging with correlation IDs

**Phase 1.3 - Observability & Error Handling**:

- `/backend/pkg/logging/logger.go` - Structured logging interface
- `/backend/pkg/errors/errors.go` - Standardized error types with context
- `/backend/pkg/middleware/correlation.go` - Correlation ID middleware
- `/backend/pkg/health/checks.go` - Health check patterns

**Phase 1.4 - Event Bus Foundation**:

- `/backend/pkg/events/bus.go` - Event bus interface
- `/backend/pkg/events/inmemory/bus.go` - In-memory event bus implementation
- `/backend/internal/events/card_events.go` - Domain event definitions
- `/backend/services/image-renderer/listener.go` - Image generation event listener

**Phase 1.5 - gRPC Service Contracts**:

- `/backend/proto/card_service.proto` - Protocol Buffer service definitions
- `/backend/services/card-generator/grpc/server.go` - gRPC server implementation
- `/backend/services/image-renderer/grpc/client.go` - gRPC client for card generator
- `/backend/pkg/registry/service.go` - Service discovery/registry

**Microservices (Phase 2)**:

- `/backend/services/card-generator/main.go` - Card Generator service
- `/backend/services/card-generator/handlers/generate.go` - Card generation endpoint
- `/backend/services/image-renderer/main.go` - Image Renderer service
- `/backend/services/importer/main.go` - Importer service

**Configuration & Deployment**:

- `/docker-compose.prod.yml` - Production configuration override
- `/DEPLOYMENT.md` - Deployment documentation
- `/scripts/deploy.sh` - Deployment automation
- `/vscode-settings.json` - Open Router configuration for VSCode

**Demo & Documentation**:

- `/demo/sample-cards.csv` - Sample card data for demo
- `/demo/demo-script.sh` - Automated demo workflow
- `/demo/README.md` - Demo walkthrough

**Modified Files**:

- `/frontend/src/server/api/routers/` - Update to use API Gateway with correlation IDs
- `/backend/Makefile` - Add `proto-generate` target for Protocol Buffers
- `/backend/go.mod` - Add grpc, protobuf, chi, and logging dependencies
- `/docker-compose.yml` - Update service ports for gRPC (5051, 5052, etc.)

---

## Notes for Future Agents

- **Architecture Philosophy**: This project uses clean architecture with domain-driven design. Maintain separation of concerns when making changes.
- **Service Communication**: Services communicate via HTTP with docker DNS. Keep payloads small and responses fast.
- **Database**: PostgreSQL schema is source of truth. Migrations should be idempotent and version controlled.
- **Testing**: Maintain unit/integration/e2e test coverage. Add service integration tests for cross-service communication.
- **Quality Standards**: Code should pass: `go fmt`, `go vet`, tests >80% coverage, and architecture compliance (SOLID/Clean/DDD principles).
- **Open Router Integration**: Used for enhanced agentic development workflow, not for card generation (keep that deterministic/controlled).

**Remember**: This modernization is about enabling scalability and maintainability while preserving the solid business logic already built. Focus on clean interfaces and service boundaries.
