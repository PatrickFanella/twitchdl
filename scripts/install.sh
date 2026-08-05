#!/usr/bin/env bash

set -Eeuo pipefail

repo_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd -P)"
prefix="${HOME:?HOME is not set}/.local"
binary_dir="$prefix/bin"
completion_dir="$prefix/share/bash-completion/completions"

mkdir -p -- "$binary_dir" "$completion_dir"
cd -- "$repo_dir"
go build -trimpath -o "$binary_dir/twitchdl" ./cmd/twitchdl
"$binary_dir/twitchdl" completion bash >"$completion_dir/twitchdl"

printf 'Installed %s\n' "$binary_dir/twitchdl"
printf 'Installed %s\n' "$completion_dir/twitchdl"
