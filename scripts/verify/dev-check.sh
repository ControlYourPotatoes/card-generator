#!/usr/bin/env bash
set -euo pipefail

echo "Running spec validation..."
bash scripts/verify/spec-check.sh

echo "Running backend compile check..."
(cd backend && go test ./... -run TestDoesNotExist >/dev/null)

echo "Running backend unit packages..."
(cd backend && go test ./pkg/... -v)

echo "Development checks passed."
