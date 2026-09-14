#!/usr/bin/env bash
#MISE description="Run Helm lint and secret key validation"
set -euo pipefail

echo "Running helm lint..."
helm lint --strict ./charts

echo "Validating rendered secret key..."
encoded_key="$(
  helm template dependencytrack ./charts \
    --set backend.postgresql.password=test \
    --set backend.postgresql.username=test |
    awk '
      /^# Source: dependencytrack\/templates\/secret-key.yaml$/ { in_secret = 1; next }
      in_secret && /^---$/ { exit }
      in_secret && /^[[:space:]]+secret\.key:/ {
        value = $2
        gsub(/^"|"$/, "", value)
        print value
        exit
      }
    '
)"

if [[ -z "$encoded_key" ]]; then
  echo "Could not find secret.key in rendered Secret" >&2
  exit 1
fi

if [[ "$(printf '%s' "$encoded_key" | base64 --decode | wc -c | tr -d ' ')" != "32" ]]; then
  echo "Rendered secret.key must decode to exactly 32 bytes" >&2
  exit 1
fi
