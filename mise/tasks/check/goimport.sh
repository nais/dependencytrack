#!/usr/bin/env bash
#MISE description="Run goimports"
set -euo pipefail

echo "Running goimports..."
files_needing_updates="$(
  find . -type f -name '*.go' ! -path './pkg/dependencytrack/client/*' \
    -exec go tool golang.org/x/tools/cmd/goimports -l {} +
)"

if [ -n "$files_needing_updates" ]; then
  echo "$files_needing_updates"
  exit 1
fi
