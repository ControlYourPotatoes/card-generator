# Codebase Assessment (2026-02-23)

## Summary
The project has solid foundations in backend domain modeling and Docker/devcontainer setup. The main risk is mismatch between planning docs and current implementation status, plus one known backend compile blocker. Best starting point is to stabilize baseline checks, then align frontend routing with existing gateway contracts.

## Architecture Snapshot
- Backend: Go monolith with emerging microservice layout under `backend/services/`.
- Frontend: Next.js/T3 scaffold with incomplete integration to API gateway.
- Infra: `docker-compose.yml` for Postgres/Adminer/Gateway/Frontend (some services still placeholders).
- Workflow: Spec-driven guardrails now added (`specs/`, `docs/agent-workflow.md`, verification scripts).

## Implemented vs Planned vs Stale
### Implemented
- API gateway service scaffold with `/api/v1` stubs and health endpoint.
- Docker + devcontainer foundation for local development.
- Domain model and template/image assets present in backend.
- New spec-first workflow checks and templates.

### Planned (not complete)
- Full frontend card workflow wiring via tRPC + gateway.
- Real microservice handlers behind gateway stubs.
- SVG template ingestion/composition end-to-end production path.

### Stale / conflicting signals
- `PROGRESS.md` includes older inline notes that look like scratchpad logs.
- Several planning docs overlap (`plan.md`, `SVG_MIGRATION_PLAN.md`, `SVG_TEMPLATE_SYSTEM_PLAN.md`) and can conflict without a canonical status index.
- Root `TODO` is broad and architecture-heavy; needs prioritization by current phase.

## Current Technical Blockers
1. Backend compile issue:
   - `backend/internal/analysis/tagger/detector.go:7` imports `github.com/ControlYourPotatoes/card-generator/internal/analysis/types`
   - Module path should include `/backend` based on `backend/go.mod`

## Top 3 Starting Points (Ranked)
1. **Baseline Stabilization (Recommended)**
- Scope: Fix compile blocker, make `make dev-check` pass baseline.
- Why first: Ensures every next task has fast, trustworthy feedback.
- Tradeoff: Limited visible product progress in short term.

2. **Frontend-to-Gateway Integration Completion**
- Scope: Wire `cardRouter` in root router, replace starter page, validate contracts.
- Why second: Converts scaffolding into usable product loop.
- Tradeoff: Risks chasing integration bugs if baseline is unstable.

3. **Docs Consolidation + Canonical Roadmap**
- Scope: Archive outdated docs, create current roadmap index, map active phase.
- Why third: Reduces confusion and context switching.
- Tradeoff: Documentation-only progress unless paired with implementation task.

## Recommended Start-Here Choice
Start with **Baseline Stabilization**:
1. Fix known backend import path issue.
2. Re-run `make dev-check`.
3. Create a focused spec for frontend gateway integration as the first feature slice.

## Next Slice Candidate Spec
- Suggested filename: `specs/frontend-gateway-integration.md`
- Goal: Complete practical frontend wiring to existing gateway contracts with low-risk scope.
