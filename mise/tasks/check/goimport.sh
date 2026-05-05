#!/usr/bin/env bash
#MISE description="Run goimports"
set -euo pipefail

echo "Running goimports..."
find . -type f -name '*.go' ! -path './pkg/dependencytrack/client/*' \
  -exec go tool golang.org/x/tools/cmd/goimports -l -w {} +
