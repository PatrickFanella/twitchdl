# Troubleshooting

Run `twitchdl doctor` first. On Arch, install dependencies with:

```bash
sudo pacman -S --needed yt-dlp streamlink ffmpeg
```

An offline channel fails unless `--wait` or `--wait-timeout` is explicit. If a
requested quality is unavailable, retry with `--quality best`.

When finalization fails, retain the reported `.part` file. A typical manual MKV
recovery attempt is:

```bash
ffmpeg -i recording.mkv.part -map 0 -c copy recovered.mkv
ffprobe -v error recovered.mkv
```

Quote URLs containing shell metacharacters. Version 1 does not support private
or subscriber-only content, browser cookies, OAuth tokens, or Twitch API keys.
