#!/usr/bin/env bash
#MISE description="Run go fix"
set -euo pipefail

echo "Running go fix..."
go fix ./...

if ! git diff --quiet -- .; then
  echo "go fix modified tracked files. Please run 'go fix ./...' and commit the changes."
  git diff -- .
  exit 1
fi
