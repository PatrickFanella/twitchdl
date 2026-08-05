# Live recording

Live capture has two phases: Streamlink records raw segmented media, then
`ffmpeg` remuxes it without re-encoding and `ffprobe` validates it. The final
filename is published atomically only after validation.

Press Ctrl+C to stop. `twitchdl` interrupts Streamlink, preserves recorded data,
and attempts finalization with a bounded recovery context. A successfully
finalized interrupted recording is reported but the command exits `130` to
retain standard interrupt semantics.

If remuxing or validation fails, the `.part` file is retained. Do not delete it;
the error reports its path for manual recovery. MKV is the default because it is
more tolerant of stream discontinuities and interruption. MP4 is optional.

`--from-start` is best effort and is constrained by Twitch's current HLS window.
It cannot recover segments Twitch no longer exposes.
