# twitchdl

`twitchdl` downloads public Twitch clips and VODs and records public live
streams. It is a Go orchestration CLI around `yt-dlp`, Streamlink, `ffmpeg`, and
`ffprobe`, with an application layer designed for a future TUI.

> Download only media you are permitted to save. You are responsible for
> complying with copyright law and Twitch's terms.

## Status

The CLI foundation is functional. It supports direct public URLs and channel
names; discovery, scheduling, authenticated content, and Twitch API credentials
are intentionally deferred.

## Requirements

- Go 1.26 or newer to build
- `yt-dlp`
- `streamlink`
- `ffmpeg` and `ffprobe`

On Arch Linux:

```bash
sudo pacman -S --needed yt-dlp streamlink ffmpeg
```

## Quick start

```bash
make build
./bin/twitchdl doctor

# Twitch clip
./bin/twitchdl clip 'https://clips.twitch.tv/CLIP_SLUG'

# Whole or ranged VOD
./bin/twitchdl vod 'https://www.twitch.tv/videos/1234567890'
./bin/twitchdl vod 'https://www.twitch.tv/videos/1234567890' 1:20:00 1:35:00

# Live now, indefinitely until Ctrl+C
./bin/twitchdl live twitchdev

# Wait indefinitely, or for a finite window
./bin/twitchdl live twitchdev --wait
./bin/twitchdl live twitchdev --wait-timeout 6h

# Record two hours and request MP4 instead of default MKV
./bin/twitchdl live twitchdev --duration 2h --container mp4
```

Install for the current user:

```bash
./scripts/install.sh
```

## Recording safety

Live data is written to a recoverable `.part` file. On normal completion or
Ctrl+C, `ffmpeg` remuxes it without re-encoding, `ffprobe` validates the result,
and only then is the final file published. If finalization fails, the raw file
is retained and reported. MKV is the default because it tolerates interrupted
and discontinuous streams better than MP4.

## Development

```bash
make verify
```

The test suite is offline: it uses injected backends and fake process contracts,
not arbitrary Twitch downloads.

An opt-in integration check generates and serves a two-second HLS fixture
locally, then records and validates it with the real media tools:

```bash
make integration
```

## Documentation

- [Usage](docs/usage.md)
- [Live recording](docs/live-recording.md)
- [Architecture and future TUI](docs/architecture.md)
- [Troubleshooting](docs/troubleshooting.md)
- [Architectural decisions](docs/adr/README.md)

## Open Pilot

This repository keeps the base Open Pilot issue template and PR template.
Bootstrap labels after creating the repository:

```bash
open-pilot labels bootstrap subculture-collective/twitchdl
```

## License

GPL-3.0-only, preserved from the source template. See `LICENSE`.
