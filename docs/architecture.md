# Architecture

`cmd/twitchdl` owns process startup and signal context. `internal/cli` parses and
renders Cobra commands. `internal/app` owns use cases and emits UI-neutral typed
events. `internal/backend` adapts `yt-dlp`, Streamlink, `ffmpeg`, and `ffprobe`.
`internal/process` owns bounded child-process lifetime. `internal/media` owns
timestamps, Twitch input normalization, output paths, and filename safety.

The planned TUI should depend on `internal/app` and implement a new reporter. It
must not execute `twitchdl` or parse terminal strings. No TUI framework is a
dependency yet; Bubble Tea can be evaluated when interaction requirements exist.

There is no persistence layer, database, server, or Twitch API client in v1.
