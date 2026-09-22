# Strong UHD Google TV Stick — scrcpy-tui v0.3.1

Qualification performed on 2026-09-21. Android identifies this target as a
`Strong` brand, `SkyworthDigital` manufacturer, `UHD Google TV Stick` model and
`KHH` device. This is distinct from the Homatics Dongle G 4K qualified earlier.
ADB addresses and device serials are omitted.

| Item | Value |
| --- | --- |
| Release | `v0.3.1` |
| Host | Debian GNU/Linux 13, Linux amd64 |
| ADB | 1.0.41, Debian platform-tools 34.0.5 |
| scrcpy | 4.1 |
| Transport | Existing authorized ADB over Wi-Fi |
| Android | 12, SDK 31 |
| Firmware | `STTB.220726.001.C1.0.4_20241113 release-keys` |
| Security patch | 2024-09-05 |

The published amd64 and arm64 archives matched `SHA256SUMS`. The executable
extracted from the amd64 archive reported `0.3.1` and was used for the tests.
An isolated temporary profile was used throughout.

## Results

| Scenario | Result |
| --- | --- |
| ADB inventory and target selection | Pass |
| `Léger Wi-Fi`, `Très léger` and `Qualité` presets | Pass: each session reached its two-second time limit and returned cleanly |
| Quoted custom window title | Pass with `--window-title 'Strong qualification'` |
| Windowed video | Pass: OpenGL renderer, 800×450 texture, clean return after four seconds |
| Audio capture with `--require-audio` | Pass: temporary Ogg/Opus recording created during a five-second session |
| SIGINT during a running session | Pass: exit 130 |
| Disappearance before launch | Pass: ADB disconnection caused exit 1 with an unavailable-device message; reconnection restored `device` |
| Profile persistence | Pass: file mode `0600`, last preset saved |
| Invalid-profile recovery | Pass: original invalid bytes preserved in one recovery copy; valid profile saved |

The temporary audio recording establishes capture and encoding, but does not
prove audible quality. Audio, raw scrcpy logs and the target's ADB identifier
were deleted with the temporary directory after qualification.

## Limits

USB transport, keyboard or pointer input forwarding, audible sound quality,
`offline` and `unauthorized` inventory states, long-duration stability and Linux
arm64 hardware were not verified on this device. The automated CI checks and
other devices' qualification records do not extend these claims.
