# Agent Instructions

This repository is the Go `twitchdl` CLI, derived from
`subculture-collective/project-template-go-worker-postgres-docker`.

## Priorities

1. Keep the CLI runnable, recoverable, and suitable for a future TUI.
2. Preserve Open Pilot issue and PR templates.
3. Keep generated caches, secrets, and local data out of git.
4. Keep media orchestration in `internal/app`; UI layers only adapt it.
5. Never delete a partial recording after failed finalization.
6. Update this file and `README.md` when stack conventions change.

## Verification

Run:

```bash
make verify
```

Use narrower commands only when an issue explicitly asks for a smaller check.
