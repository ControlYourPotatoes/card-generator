# Spec-Driven Agent Workflow

This project uses a spec-first workflow to control scope, reduce unnecessary complexity, and keep testing practical.

## Workflow
1. Create spec from template:
   - `cp specs/TEMPLATE.md specs/<feature>.md`
2. Complete the Discussion Gate before coding.
3. Implement in small slices only.
4. Run verification commands.
5. Record evidence in spec and `PROGRESS.md`.

## Discussion Gate Rules
Before writing code, the spec must define:
- Problem and unchanged behavior.
- Out-of-scope boundaries.
- Complexity budget (files/APIs/dependencies).
- Measurable impact (baseline + target + method).

If any of these are missing, implementation is blocked.

## Testing Philosophy
- Prefer real behavior over heavy mocking.
- Unit tests for pure functions and rule logic.
- Integration/contract tests for workflows and boundaries.
- A change is not done unless acceptance checks pass with evidence.

## Definition of Done
- Spec acceptance criteria checked.
- Verification commands passed.
- Measurable impact section updated with evidence.
- `PROGRESS.md` updated with factual status and date.

## Recommended Commands
- Validate specs:
  - `make spec-check`
  - `make verify-spec SPEC=specs/<feature>.md`
- Dev guardrail checks:
  - `make dev-check`
