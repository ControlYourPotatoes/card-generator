# Docs Consolidation Plan (2026-02-23)

## Objective
Reduce documentation drift and create a single reliable reference path for planning and execution.

## Proposed Canonical Structure
- `docs/agent-workflow.md`: Development operating model (spec-first process).
- `docs/codebase-assessment.md`: Current architecture/status snapshot.
- `docs/docs-consolidation.md`: This maintenance plan.
- `specs/`: Active, feature-level specs with measurable impact.
- `PROGRESS.md`: Factual implementation status log (no scratchpad notes).

## Inventory and Actions
| File | Status | Action | Notes |
|---|---|---|---|
| `PROGRESS.md` | Active | Update | Remove scratchpad-style trailing notes; keep verified facts only |
| `plan.md` | Active reference | Update | Keep as macro roadmap; add banner pointing to `specs/` for execution |
| `SVG_TEMPLATE_SYSTEM_PLAN.md` | Active reference | Keep | Detailed SVG design reference |
| `SVG_MIGRATION_PLAN.md` | Active reference | Keep | Migration strategy reference |
| `TODO` | Legacy/mixed | Update | Convert into prioritized backlog linked to specs |
| `promptcreator.md` | Active tool | Update | Add current handoff prompt template with references |
| `.devcontainer/README.md` | Active | Keep | Clear runtime/dev guidance |
| `container_setup.md` | Legacy support note | Archive | Move to `docs/archive/` unless still actively used |
| `backend/README.md` | Active | Keep | Backend usage reference |
| `frontend/README.md` | Boilerplate | Update | Replace T3 starter content with project-specific usage |

## Conflicts to Resolve
1. `PROGRESS.md` currently mixes status and internal work notes.
2. Multiple high-level plans can imply different current priorities.
3. Frontend docs still reflect template defaults, not current architecture.

## Consolidation Execution Plan
1. Clean `PROGRESS.md` to verified status entries only.
2. Add a short “Docs Index” section to `plan.md` that points to canonical docs.
3. Move stale docs to `docs/archive/` with a one-line reason in each moved file header.
4. Replace generic frontend readme content with project-specific run/integration instructions.
5. Require each new feature to have a `specs/*.md` entry before implementation.

## Guardrails for Future Docs
- No roadmap changes without date + owner + validation evidence.
- No feature execution plans outside `specs/`.
- If a doc becomes outdated, mark it with `Status: Archived` and link replacement.
