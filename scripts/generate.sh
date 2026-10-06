#!/usr/bin/env bash
# Generate Go server stubs and the SvelteKit TypeScript client from spec/openapi.yaml.
# Not wired in this slice (no oapi-codegen / openapi-typescript yet).
set -euo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"
echo "OpenAPI codegen is not set up yet. Spec: ${root}/spec/openapi.yaml" >&2
exit 1
