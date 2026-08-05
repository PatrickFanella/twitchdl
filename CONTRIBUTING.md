# Contributing to twitchdl

Keep UI code thin: media orchestration belongs in `internal/app`, external tool
adapters in `internal/backend`, and terminal concerns in `internal/cli`.

Before opening a pull request, run:

```bash
make verify
```

Tests must not contact Twitch by default. Use fake runners or local media
fixtures, preserve partial recordings on failure, and never commit downloaded
media, cookies, tokens, or API secrets.

## Development flow

1. Create or pick an issue before making non-trivial changes.
2. Create a branch from `main`.
3. Make the smallest change that satisfies the issue.
4. Run verification:

   ```bash
   make verify
   ```

5. Open a pull request using the PR template.
6. Link the issue with `Closes #N` when appropriate.

## Commit style

Prefer concise conventional-style subjects:

```text
feat: add recording quality selection
docs: update Open Pilot guidance
fix: retain partial recording after remux failure
chore: update media dependencies
```

## Open Pilot issues

Use `.gitea/ISSUE_TEMPLATE/open-pilot-task.yaml` for tasks that Open Pilot should implement.

Open Pilot requires:

- `ready-for-agent`
- `agent:queued`
- a deterministic `Test Command`

The issue form adds `ready-for-agent`; add `agent:queued` only after the issue is complete and ready for automation.

Good Open Pilot tasks are:

- one focused outcome
- clear context and constraints
- explicit out-of-scope notes
- verifiable from the repository root

## Verification commands

Use the narrowest deterministic command that proves the change.

Default full check:

```bash
make verify
```

Examples:

```bash
go test ./...
go test -race ./...
make integration
```

## Project boundaries

Keep Twitch API discovery, authenticated media, scheduling, persistence, and TUI
framework dependencies out of v1 unless an approved issue explicitly expands
scope and defines credential or lifecycle safety.
