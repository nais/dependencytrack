#!/usr/bin/env bash
#MISE description="Run gosec"
set -euo pipefail

echo "Running gosec..."
go tool github.com/securego/gosec/v2/cmd/gosec --exclude-generated -terse --exclude-dir=pkg/dependencytrack/client ./...
