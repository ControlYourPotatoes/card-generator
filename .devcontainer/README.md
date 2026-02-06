# Dev Container vs Docker Compose — Separation of Concerns

## Architecture

```
┌──────────────────────────────────────────────────┐
│  Dev Container (.devcontainer/)                  │
│  PURPOSE: Coding environment                     │
│  ─────────────────────────────────               │
│  • Go 1.23 + gopls + dlv + golangci-lint         │
│  • Node.js 20 + npm + Prisma CLI                 │
│  • Docker CLI (talks to host Docker daemon)      │
│  • VS Code extensions (Go, ESLint, Prettier...)  │
│  • You write code and run services HERE          │
│                                                  │
│  go run ./backend/services/api-gateway/          │
│  cd frontend && npm run dev                      │
└──────────────────────────────────────────────────┘
         │ uses host Docker daemon
         ▼
┌──────────────────────────────────────────────────┐
│  Docker Compose (docker-compose.yml)             │
│  PURPOSE: Infrastructure & deployment            │
│  ─────────────────────────────────               │
│  • postgres:16 (port 5433)                       │
│  • adminer (port 8081)                           │
│  • api-gateway (port 8080) ← for deployment      │
│  • frontend (port 3000)    ← for deployment      │
└──────────────────────────────────────────────────┘
```

## Development Workflow

### Option A: Dev Container + Compose Infrastructure (Recommended)

1. Open workspace in VS Code → "Reopen in Container"
2. Start infrastructure: `docker compose --env-file .env.docker up -d postgres adminer`
3. Run backend directly: `cd backend && go run ./services/api-gateway/`
4. Run frontend directly: `cd frontend && npm run dev`

**Why?** Hot reload, debugger attach, fast iteration. No rebuild on every change.

### Option B: Full Docker Compose (Deployment Testing)

1. From host terminal: `docker compose --env-file .env.docker up -d`
2. All services run in containers (postgres, api-gateway, frontend, adminer)

**Why?** Test the full containerized stack as it would run in production.

## Key Design Decisions

| Concern          | Dev Container                 | Docker Compose                  |
| ---------------- | ----------------------------- | ------------------------------- |
| **When to use**  | Daily coding                  | Integration testing, deployment |
| **Go backend**   | `go run` directly             | Built into Docker image         |
| **Frontend**     | `npm run dev` with hot reload | Built into Docker image         |
| **PostgreSQL**   | Runs via compose (port 5433)  | Runs via compose (port 5433)    |
| **Code changes** | Instant (no rebuild)          | Requires `docker compose build` |
| **Debugging**    | Full debugger support         | Limited (container logs)        |

## Port Map

| Port | Service          | Used By                                |
| ---- | ---------------- | -------------------------------------- |
| 3000 | Next.js frontend | Dev container (npm run dev) OR compose |
| 5433 | PostgreSQL       | Always via compose                     |
| 8080 | API Gateway      | Dev container (go run) OR compose      |
| 8081 | Adminer          | Always via compose                     |

## Environment Variables

The dev container sets `DATABASE_URL` pointing to `host.docker.internal:5433` so the Go backend and Prisma can reach the PostgreSQL container running via docker-compose.

If running **outside** the dev container (directly on host), use `localhost:5433` instead.
