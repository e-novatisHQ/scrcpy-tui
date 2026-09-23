# WIN-08 — Windows 11 qualification protocol

Status: **NOT RUN**. This is a protocol, not a qualification record. Native hosted
Windows CI and synthetic MSI upgrades do not replace this desktop/device campaign.
Use a clean Windows 11 x64 account or VM snapshot and an authorized Android device.
Do not use a production installation. Record operator, UTC date, Windows build,
terminal/PowerShell versions, adb/scrcpy versions, Android model/firmware, release
commit, MSI/ZIP SHA-256 and test transport. Replace serials, usernames and addresses
with consistent synthetic identifiers in shared evidence. No credentials or pairing
codes belong in Git, Plane or CI artifacts.

| Step | Action | Required observation | Result |
|---|---|---|---|
| 1 | Record baseline user/system PATH, installed apps and empty target directory | No preexisting scrcpy-tui MSI, clean snapshot available | Not run |
| 2 | Check MSI hash, provenance and signature status | Exact candidate identified; unsigned status explicitly recorded if applicable | Not run |
| 3 | Double-click MSI as a standard user | Installation completes without manual file moves or elevation | Not run |
| 4 | Open a fresh Windows Terminal from Explorer | `scrcpy-tui --version` resolves expected binary; other PATH entries intact | Not run |
| 5 | Install prerequisites following Windows guide | adb and scrcpy version/source recorded; missing dependency diagnostics verified in an isolated PATH | Not run |
| 6 | Authorize USB debugging on Android | Device state `device`; unauthorized/offline states handled; no unintended target | Not run |
| 7 | Preview and launch every built-in preset over USB | Matching target, visible video, audio behavior and input explicitly recorded | Not run |
| 8 | Close scrcpy; repeat with Ctrl+C | Menu/terminal restored; no owned descendant remains; preexisting adb survives | Not run |
| 9 | Establish authorized ADB Wi-Fi connection separately; repeat 7–8 | Real network session succeeds; video/audio/input and cleanup recorded | Not run |
| 10 | Save Unicode/space-containing preset, restart and recover invalid temporary profile | Persistence and one exact backup; user ACL protected | Not run |
| 11 | Upgrade from an earlier test MSI through double-click | One product registration/PATH entry, new binary version, existing profile unchanged | Not run |
| 12 | Uninstall in Windows Settings; open a fresh terminal | Installed files/PATH removed, other PATH entries and user profiles intact | Not run |
| 13 | Reinstall candidate and exercise ZIP separately | Both documented launch paths work; no stale binary mistaken for the candidate | Not run |

On failure, stop the affected scenario, preserve sanitized logs and restore the
VM snapshot or uninstall only the test product. Do not kill a shared adb server.
A failure or unexecuted row prevents WIN-08 closure. Include screenshots only after
redaction; a screenshot of a process list alone cannot prove actual video/audio.
Record pass/fail per row and links to evidence in a new dated result document.
Acceptance requires all applicable rows to pass, no blocking defect and an explicit
operator verdict. WIN-09 signature verification is separate and cannot be inferred
from SHA-256 or GitHub attestations.
