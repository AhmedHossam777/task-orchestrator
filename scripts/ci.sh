#!/usr/bin/env bash
set -euo pipefail

echo "==> Running CI checks for task-orchestrator"

echo "==> Go version"
go version

echo "==> Building"
go build ./...

echo "==> Vetting"
go vet ./...

echo "==> Testing"
go test ./...

echo "==> CI checks passed"
