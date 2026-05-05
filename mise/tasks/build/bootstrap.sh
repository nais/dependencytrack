#!/usr/bin/env bash
#MISE description="Build the bootstrap binary"
set -euo pipefail

go build -o bin/bootstrap ./cmd/bootstrap/...
