# Xiaomi Mi A3 — scrcpy-tui v0.3.1

Qualification performed on 2026-09-21 with the executable extracted from the
published Linux amd64 archive. Android identifies the phone as Xiaomi Mi A3,
device `laurel_sprout`, Android 11 (SDK 30), build
`qssi-user 11 RKQ1.200903.002 V12.0.26.0.RFQMIXM`, security patch 2022-08-01.
It was connected through authorized ADB USB. The host was Debian GNU/Linux 13
amd64 with ADB 1.0.41 (platform-tools 34.0.5) and scrcpy 4.1. The ADB serial
and raw logs are omitted.

| Scenario | Result |
| --- | --- |
| ADB inventory and target selection | Pass |
| `Léger Wi-Fi`, `Très léger` and `Qualité` presets | Pass: each two-second session returned cleanly |
| Quoted `--window-title` argument | Pass |
| Windowed video | Pass: OpenGL renderer, 800×368 texture, clean return after four seconds |
| SIGINT during a running session | Pass: exit 130 |
| Audio capture with `--require-audio` | **Failed**: scrcpy reported an audio-capture error; the 167-byte Ogg/Opus file held only headers, no audio data, while the process still exited 0 |

The empty recording is not counted as an audio success. The failure may depend
on device or Android 11 state; its cause was not established in this run. The
recording and raw logs were retained only in a private temporary directory and
deleted after inspection.

USB video and preset launching are qualified for this combination. Audio
forwarding, keyboard and pointer control, real `offline` and `unauthorized`
states, long-duration behavior and other Mi A3 firmware versions remain
unqualified.
