#!/bin/bash

set -euo pipefail

# LaunchAgent-friendly entrypoint.
# Prefer a pre-built binary for faster startup and fewer deps.
if [ -x "./antigravity-proxy" ]; then
  exec ./antigravity-proxy
fi

# Fallback for dev setups.
if command -v mise >/dev/null 2>&1; then
  exec mise run run
fi

echo "Error: missing ./antigravity-proxy binary and mise is not installed" >&2
exit 1
