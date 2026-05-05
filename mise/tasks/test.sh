#!/usr/bin/env bash
#MISE description="Run tests"
set -euo pipefail

IGNORED_PATH="pkg/dependencytrack/client"
GO_PACKAGES=$(go list ./... | grep -v "$IGNORED_PATH")

# shellcheck disable=SC2086
go test -v --race --cover --coverprofile=cover.out $GO_PACKAGES
