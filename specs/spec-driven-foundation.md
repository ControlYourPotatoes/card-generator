# Spec: spec-driven-foundation

## 1. Problem
- Problem statement: Development scope expands too easily, complexity creeps in, and testing can become overly mocked and low-value.
- Users affected: Primary maintainer of this project.
- Why now: Project is being resumed and needs a reliable operating model before major implementation work.

## 2. Scope
- In scope:
  - Add workflow guardrails for spec-first development.
  - Enforce measurable impact fields per feature spec.
  - Add a discussion gate before implementation.
  - Add practical verification commands.
- Out of scope:
  - Implementing product features in card generation pipeline.
  - Full CI/CD redesign.
- Assumptions:
  - Development happens mostly inside devcontainer.

## 3. Discussion Gate (Required Before Coding)
- Decision 1: Problem clarity
  - What exact behavior is changing? Every implementation task starts with a spec and explicit acceptance criteria.
  - What stays unchanged? Existing project architecture and phase plans.
- Decision 2: Complexity budget
  - Max files touched: 8
  - Max new public APIs: 0
  - Max new dependencies: 0
  - Fallback if budget is exceeded: Split into a follow-up spec.
- Decision 3: Evidence of value
  - What will be measured: Spec completeness, repeatable verification usage.
  - How we will capture baseline: Current repo has plans but no enforced spec checks.
  - How we will validate improvement: `make spec-check` and `make dev-check` become standard pre-merge commands.

## 4. Interfaces and Contracts
- API/CLI/DB contracts affected: None.
- Backward compatibility plan: No runtime behavior changes.
- Migration notes: Team starts creating specs in `specs/` for each feature.

## 5. Acceptance Criteria
- [x] Spec template exists with measurable impact section.
- [x] Workflow doc exists with discussion gate.
- [x] Cursor coding guardrails are always-on and include architecture preferences.
- [x] Make targets exist for spec validation and dev checks.
- [x] Validation scripts exist and are executable.

## 6. Measurable Impact
| Metric | Baseline | Target | Measurement Method |
|---|---:|---:|---|
| Features with explicit acceptance criteria | 0% | 100% | Count features with corresponding spec files |
| Features with measurable impact fields | 0% | 100% | Run `make spec-check` |
| Pre-implementation discussion gate usage | 0% | 100% | Confirm section completion in each spec |

## 7. Practical Test Plan
- Unit tests:
  - Not required for docs/scripts-only change.
- Integration tests:
  - Validate scripts run in repo context.
- Contract tests:
  - Not applicable.
- Manual verification:
  - `make spec-check`
  - `make verify-spec SPEC=specs/spec-driven-foundation.md`

## 8. Rollout and Rollback
- Rollout steps:
  - Adopt spec-first process for next feature.
  - Run `make spec-check` before coding.
- Rollback trigger:
  - Workflow blocks useful work with no value.
- Rollback steps:
  - Revert new scripts/targets and keep docs only.

## 9. Implementation Plan (Small Slices)
1. Slice 1:
   - Scope: Create template and workflow documentation.
   - Acceptance: Files present and reviewable.
2. Slice 2:
   - Scope: Add script and Makefile checks.
   - Acceptance: Commands run successfully.

## 10. Evidence Log
- Date: 2026-02-23
- Commands run:
  - `make spec-check`
  - `make verify-spec SPEC=specs/spec-driven-foundation.md`
- Results:
  - Pending run after file creation.
- Risks left:
  - `dev-check` can fail while existing backend compile issues remain unresolved.
