# Homatics Dongle G 4K — scrcpy-tui v0.3.1

Qualification performed on 2026-09-21. The target was initially described as a
Strong dongle, but Android identifies it as `Homatics Dongle G 4K` by
`SEI Robotics`. A separate connected device reports the Strong brand. This
record is therefore attributed to Homatics; the Strong-branded device was not
qualified in this run. ADB addresses and device serials are omitted.

## Environment

| Item | Value |
| --- | --- |
| Release | `v0.3.1` |
| Release archive | `scrcpy-tui_0.3.1_linux_amd64.tar.gz` |
| Archive SHA-256 | `9b234c4d095b04f620ee986d3eb8b252a254b5f6d00661a7bf3abee983b94ae8` |
| Host | Debian GNU/Linux 13, Linux amd64 |
| ADB | 1.0.41, Debian platform-tools 34.0.5 |
| scrcpy | 4.1 |
| Transport | Existing authorized ADB over Wi-Fi |
| Manufacturer / brand | SEI Robotics / Homatics |
| Model / device | Dongle G 4K / YQX |
| Android | 14, SDK 34 |
| Firmware | `UKG3.250826.001.6935` |
| Security patch | 2025-11-05 |

The released archive checksum matched `SHA256SUMS`, its GitHub artifact
attestation matched the archive digest, and the extracted executable reported
`0.3.1`.

## Results

| Scenario | Result | Evidence |
| --- | --- | --- |
| Real ADB inventory | Pass | Eight authorized devices detected; target selected from the launcher's previously saved device, without recording its address |
| `Léger Wi-Fi` preset | Pass | Server started and reached the two-second time limit cleanly |
| `Très léger` preset | Pass | Server started and reached the two-second time limit cleanly |
| `Qualité` preset | Pass | Server started and reached the two-second time limit cleanly |
| Quoted custom argument | Pass | `--window-title 'Qualification Strong'` was preserved during a real launch |
| Windowed video session | Pass | Android 14 device recognized; OpenGL renderer created an 800×450 texture and returned cleanly after five seconds |
| SIGINT / Ctrl+C contract | Pass | Exit 130; no matching scrcpy process remained |
| SIGTERM contract | Pass | Exit 130; no matching scrcpy process remained |
| scrcpy failure | Pass | Unsupported codec produced a clear diagnostic and exit 1 |
| Device disappears before launch | Pass | Target was disconnected, launch was refused as unavailable, then ADB reconnection returned to `device` |
| Profile creation | Pass | Temporary profile created with mode `0600` and last selections persisted |
| Invalid-profile recovery | Pass | Original invalid bytes preserved in one recovery file; valid profile restored with mode `0600` |
| TUI 60×18 | Pass | Device and preset surfaces rendered; `q` returned 0 |
| TUI 100×30 | Pass | Device and preset surfaces rendered; `q` returned 0 |
| TUI 160×45 | Pass | Device and preset surfaces rendered; `q` returned 0 |

All launch tests used the executable extracted from the published release. Raw
scrcpy output was kept only in a temporary directory because scrcpy prints the
full ADB inventory. The durable evidence contains no address, serial or user
profile.

## Qualification scope

This firmware combination is qualified for ADB over Wi-Fi, the three built-in
presets, video session establishment, lifecycle handling and profile recovery.

The following claims require separate evidence and are not made by this run:

- USB transport on this dongle;
- offline and unauthorized inventory states;
- audio forwarding quality;
- manual validation of keyboard, pointer or remote-control input forwarding;
- long-duration stability, latency or thermal behavior;
- other Homatics or SEI firmware versions, and the separate Strong-branded device.
