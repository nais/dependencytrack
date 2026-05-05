#!/usr/bin/env bash
#MISE description="Run go fix"
set -euo pipefail

echo "Running go fix..."
go fix ./...
