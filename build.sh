#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$script_dir"

if ! command -v go >/dev/null 2>&1; then
  echo "Go is not installed or not in PATH."
  exit 1
fi

go build -o fyne-ssh ./cmd/fyne-ssh
./fyne-ssh
