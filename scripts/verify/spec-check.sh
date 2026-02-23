#!/usr/bin/env bash
set -euo pipefail

required_sections=(
  "^## 1\\. Problem"
  "^## 2\\. Scope"
  "^## 3\\. Discussion Gate"
  "^## 5\\. Acceptance Criteria"
  "^## 6\\. Measurable Impact"
  "^## 7\\. Practical Test Plan"
  "^## 8\\. Rollout and Rollback"
  "^## 10\\. Evidence Log"
)

check_file() {
  local file="$1"
  local missing=0

  for section in "${required_sections[@]}"; do
    if ! rg -q "${section}" "${file}"; then
      echo "Missing section in ${file}: ${section}"
      missing=1
    fi
  done

  if ! rg -q "^\| Metric \| Baseline \| Target \| Measurement Method \|" "${file}"; then
    echo "Missing measurable impact table header in ${file}"
    missing=1
  fi

  if [[ "${missing}" -ne 0 ]]; then
    return 1
  fi

  echo "Spec check passed: ${file}"
}

if [[ $# -gt 0 ]]; then
  check_file "$1"
  exit 0
fi

mapfile -t specs < <(find specs -maxdepth 1 -type f -name "*.md" ! -name "TEMPLATE.md" | sort)

if [[ ${#specs[@]} -eq 0 ]]; then
  echo "No spec files found in specs/ (excluding TEMPLATE.md)"
  exit 1
fi

for spec in "${specs[@]}"; do
  check_file "${spec}"
done

echo "All specs passed validation."
