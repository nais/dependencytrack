#!/usr/bin/env bash
#MISE description="Run helm lint"
set -euo pipefail

echo "Running helm lint..."
helm lint --strict ./charts
helm lint --strict ./charts-v5 --values ./charts-v5/test-values.yaml

echo "Rendering Dependency-Track v5 wrapper chart..."
helm template dependencytrack-v5 ./charts-v5 --values ./charts-v5/test-values.yaml >/dev/null
