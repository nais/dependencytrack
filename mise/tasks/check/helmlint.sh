#!/usr/bin/env bash
#MISE description="Run helm lint"
set -euo pipefail

echo "Running helm lint..."
helm lint --strict ./charts
