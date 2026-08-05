# ADR 0001: Bootstrap twitchdl from the Go worker template

Status: accepted

## Context

The source repository was a generic Go worker with PostgreSQL, Docker Compose,
migrations, and Open Pilot conventions. `twitchdl` is a local media CLI with a
planned TUI and no server-side persistence requirement.

## Decision

Preserve the template's Git lineage, GPL-3.0-only license, Gitea conventions,
documentation shape, and `make verify` contract. Remove the worker, PostgreSQL,
migrations, Docker, Compose, and environment configuration.

Go owns orchestration, validation, process lifecycle, recoverable output, and
UI-neutral events. Mature external tools retain media responsibilities:
`yt-dlp` for clips/VODs, Streamlink for live capture, and `ffmpeg`/`ffprobe` for
finalization. Live recordings default to MKV. Public media is the only v1 scope;
Twitch API discovery and authenticated playback are deferred.

The CLI calls `internal/app`; a future TUI must call that same layer rather than
executing or parsing the CLI.

## Consequences

The executable remains small and maintainable but requires the external tools.
There is no database or daemon. A later TUI can replace terminal rendering while
reusing requests, events, backends, cancellation, and output safety.
