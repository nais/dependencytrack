#!/usr/bin/env bash
#MISE description="Generate Go client from OpenAPI spec"
set -euo pipefail

echo "Generating Go code from the OpenAPI specification..."
docker run --rm \
  -v "$(pwd):/local" \
  --platform linux/amd64 \
  openapitools/openapi-generator-cli generate \
  -i /local/schema/dtrack.json \
  -g go \
  -o /local/pkg/dependencytrack/client \
  --global-property apiTests=false,modelTests=false \
  --package-name client \
  --additional-properties=withGoMod=false \
  --additional-properties=generateInterfaces=true \
  --additional-properties=packageName=client || {
    echo "Error: docker is not running or openapi-generator-cli image failed to execute."
    exit 1
  }
