#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
command -v go >/dev/null || { echo "Установи Go: sudo apt install golang-go"; exit 1; }
go build -trimpath -ldflags="-s -w" -o goInstaller .
echo "Готово: ./goInstaller"
