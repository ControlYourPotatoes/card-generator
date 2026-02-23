# Card Generator Frontend

This frontend is a Next.js + tRPC app used to interact with the card-generator backend services through the API gateway.

## Current Architecture Role
- UI and app routes: Next.js App Router
- Type-safe server calls: tRPC routers under `src/server/api/routers`
- Gateway integration point: `src/server/api/gateway.ts`
- Contracts source of truth: `../api/contracts.ts`

## Development Modes

### Recommended: Devcontainer + Compose Infrastructure
1. Start infrastructure:
```bash
docker compose --env-file .env.docker up -d postgres adminer
```
2. Run API gateway directly:
```bash
cd backend && go run ./services/api-gateway/
```
3. Run frontend dev server:
```bash
cd frontend && npm run dev
```

### Full Compose (deployment-style check)
```bash
docker compose --env-file .env.docker up -d
```

## Local URLs
- Frontend: `http://localhost:3000`
- API Gateway: `http://localhost:8080/api/v1`
- Adminer: `http://localhost:8081`
- Postgres host port: `5433`

## Environment Notes
- Frontend expects `NEXT_PUBLIC_API_URL` to point to the gateway base path.
- In devcontainer, environment values are set in `.devcontainer/devcontainer.json`.
- For compose frontend service, `NEXT_PUBLIC_API_URL` is configured in `docker-compose.yml`.

## Verification Commands
From repo root:
```bash
make spec-check
make dev-check
```

## Related Docs
- Workflow: `docs/agent-workflow.md`
- Status snapshot: `docs/codebase-assessment.md`
- Docs consolidation: `docs/docs-consolidation.md`
