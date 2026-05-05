#!/usr/bin/env bash
#MISE description="Run deadcode"
set -euo pipefail

echo "Running deadcode..."
go tool golang.org/x/tools/cmd/deadcode -filter "pkg/dependencytrack/client" -test ./...
