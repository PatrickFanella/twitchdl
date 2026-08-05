# Usage

## Clips and VODs

```text
twitchdl clip URL [START [END]]
twitchdl vod URL [START [END]]
```

Timestamps accept seconds, `MM:SS`, or `HH:MM:SS`, with optional milliseconds.
An omitted end means the rest of the media. The end must be later than the
start.

```bash
twitchdl clip 'https://clips.twitch.tv/CLIP_SLUG'
twitchdl clip 'https://clips.twitch.tv/CLIP_SLUG' 5 20
twitchdl vod 'https://www.twitch.tv/videos/1234567890'
twitchdl vod 'https://www.twitch.tv/videos/1234567890' 1:20:00 1:35:00
```

Use `--output-dir` to choose a destination and `--force` to permit replacement.

## Live

```text
twitchdl live CHANNEL_OR_URL [options]
```

Without a wait option, an offline channel fails promptly. `--wait` waits
indefinitely; `--wait-timeout 6h` bounds the wait. They are mutually exclusive.
`--duration` starts after a stream is found. `--from-start` requests the earliest
segment still retained by Twitch, not necessarily the true broadcast beginning.

MKV is the default. Use `--container mp4` explicitly for MP4. Use `--keep-raw`
to retain Streamlink's raw capture after successful finalization.

## Diagnostics and completion

```bash
twitchdl doctor
twitchdl completion bash
twitchdl completion zsh
twitchdl completion fish
```
