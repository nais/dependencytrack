#!/usr/bin/env bash
#MISE description="Run go vet"
set -euo pipefail

IGNORED_PATH="pkg/dependencytrack/client"
GO_PACKAGES=$(go list ./... | grep -v "$IGNORED_PATH")

echo "Running go vet..."
# shellcheck disable=SC2086
go vet $GO_PACKAGES
