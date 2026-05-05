#!/usr/bin/env bash
#MISE description="Generate mocks"
set -euo pipefail

find pkg -type f -name "mock_*.go" -delete
go tool github.com/vektra/mockery/v3 --config ./mockery.yaml
find pkg -type f -name "mock_*.go" -exec go tool mvdan.cc/gofumpt -w {} \;
