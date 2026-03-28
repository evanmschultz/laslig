#!/bin/zsh

set -euo pipefail

cd "$(dirname "$0")/../.."

for d in section notice record kv list table panel paragraph statusline markdown codeblock logblock gotestout magecheck; do
  clear
  printf 'Läslig demo: %s\n\n' "$d"
  go run "./examples/$d" --format human --style always
  sleep 1
done
