# Pixel 6 (USB) and Xiaomi MiTV-MOOR2 (ADB Wi-Fi) — v0.3.1

Qualification performed on 2026-09-21 using the executable extracted from the
published Linux amd64 archive. Its `--version` returned `0.3.1`; both archives
matched the published `SHA256SUMS`. The host was Debian GNU/Linux 13 on amd64,
with ADB 1.0.41 (platform-tools 34.0.5) and scrcpy 4.1.

| Target | Connection | Android | Build | Security patch |
| --- | --- | --- | --- | --- |
| Google Pixel 6 | Authorized ADB USB | 16 | `CP1A.260405.005` | 2026-04-05 |
| Xiaomi MiTV-MOOR2 | Authorized ADB Wi-Fi | 11 | `RTM5.220609.003.2674` | 2026-02-01 |

ADB addresses and device serials are omitted. A second USB phone mentioned by
the operator was absent from `adb devices` during this run, so it was not
qualified or reset.

## Results

| Scenario | Pixel 6 USB | Xiaomi TV Wi-Fi |
| --- | --- | --- |
| ADB inventory and target selection | Pass | Pass |
| `Léger Wi-Fi` preset, two-second session | Pass | Pass |
| `Très léger` preset, two-second session | Pass | Pass |
| `Qualité` preset, two-second session | Pass | Pass |
| Quoted `--window-title` argument | Pass | Pass |
| Windowed video, four-second session | OpenGL, 360×800, clean return | OpenGL, 800×450, clean return |
| SIGINT during a running session | Exit 130 | Exit 130 |
| Five-second audio capture with `--require-audio` | Ogg/Opus stereo, 48 kHz | Ogg/Opus stereo, 48 kHz |

The first audio attempt used `--no-audio-playback` without recording; scrcpy
disabled audio because there was no playback or recording sink. It was discarded
as evidence. The successful attempt recorded a temporary Opus stream without
playing sound on the workstation. It demonstrates audio capture and encoding,
not audible quality or synchronized playback. Audio files and raw scrcpy logs
were kept only in a private temporary directory during qualification and deleted
afterward.

## Limits

No keyboard, pointer or remote-control input was sent to either target; input
forwarding remains unverified. Neither an `offline` nor an `unauthorized` ADB
state was available in the real inventory. The second USB phone and Linux arm64
host were unavailable. Long-duration stability and audio listening quality were
not evaluated. Profile recovery and multiple terminal sizes are covered by the
separate Homatics qualification, not by this target-specific run.
