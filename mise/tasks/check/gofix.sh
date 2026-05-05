#!/usr/bin/env bash
#MISE description="Run go fix"
set -euo pipefail

echo "Running go fix..."
go fix ./...

if ! git diff --exit-code; then
  echo "go fix modified files — please run 'go fix ./...' locally and commit the changes" >&2
  exit 1
fi
