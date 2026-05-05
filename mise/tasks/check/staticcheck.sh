#!/usr/bin/env bash
#MISE description="Run staticcheck"
set -euo pipefail

IGNORED_PATH="pkg/dependencytrack/client"
GO_PACKAGES=$(go list ./... | grep -v "$IGNORED_PATH")

echo "Running staticcheck..."
# shellcheck disable=SC2086
go tool honnef.co/go/tools/cmd/staticcheck -f=stylish $GO_PACKAGES
