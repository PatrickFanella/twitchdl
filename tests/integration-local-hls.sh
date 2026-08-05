#!/usr/bin/env bash

set -Eeuo pipefail

for command_name in ffmpeg ffprobe streamlink python; do
  command -v "$command_name" >/dev/null 2>&1 || {
    printf 'missing integration dependency: %s\n' "$command_name" >&2
    exit 2
  }
done

integration_tmp="$(mktemp -d)"
server_pid=''

cleanup() {
  if [[ -n "$server_pid" ]]; then
    kill "$server_pid" 2>/dev/null || true
    wait "$server_pid" 2>/dev/null || true
  fi
  case "$integration_tmp" in
    /tmp/*) rm -rf -- "$integration_tmp" ;;
    *) printf 'refusing to clean unexpected path: %s\n' "$integration_tmp" >&2 ;;
  esac
}
trap cleanup EXIT

ffmpeg -hide_banner -loglevel error \
  -f lavfi -i 'testsrc2=size=160x90:rate=10' \
  -f lavfi -i 'sine=frequency=440:sample_rate=48000' \
  -t 2 \
  -c:v libx264 -preset ultrafast \
  -c:a aac \
  -f hls -hls_time 0.5 -hls_list_size 0 \
  "$integration_tmp/index.m3u8"

python - "$integration_tmp" "$integration_tmp/port" <<'PY' &
import functools
import http.server
import pathlib
import socketserver
import sys

root = sys.argv[1]
port_file = pathlib.Path(sys.argv[2])
handler = functools.partial(http.server.SimpleHTTPRequestHandler, directory=root)
with socketserver.TCPServer(("127.0.0.1", 0), handler) as server:
    port_file.write_text(str(server.server_address[1]), encoding="utf-8")
    server.serve_forever()
PY
server_pid=$!

for _ in {1..100}; do
  [[ -s "$integration_tmp/port" ]] && break
  sleep 0.02
done
[[ -s "$integration_tmp/port" ]] || {
  printf 'local HLS server did not start\n' >&2
  exit 1
}

port="$(<"$integration_tmp/port")"
raw="$integration_tmp/recording.part"
final="$integration_tmp/recording.mkv"

streamlink \
  --output "$raw" \
  --force \
  --progress no \
  "hls://http://127.0.0.1:${port}/index.m3u8" best

ffmpeg -hide_banner -loglevel error -i "$raw" -map 0 -c copy "$final"
duration="$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "$final")"
[[ -n "$duration" ]] || {
  printf 'ffprobe returned no duration\n' >&2
  exit 1
}

printf 'local HLS integration passed: duration=%ss\n' "$duration"
