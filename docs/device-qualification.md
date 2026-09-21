# Real-device qualification

Automated tests use isolated fake adb and scrcpy executables. A release may be
promoted as tested on real hardware only when the following evidence is recorded
for each qualified combination.

Record the release tag, Linux distribution and architecture, terminal, adb and
scrcpy versions, Android device model and version (with serials removed), USB or
Wi-Fi transport, commands exercised, result, known limitations and tester date.

Exercise at least:

- inventory with usable, offline and unauthorized states;
- launch and clean return over USB and ADB Wi-Fi;
- each built-in preset and one quoted custom argument;
- Ctrl+C, SIGTERM, scrcpy failure and device disappearance before launch;
- profile creation, invalid-profile recovery and restoration;
- narrow, standard and large terminal sizes;
- archive checksum verification and `--version` from the released binary.

Store only sanitized logs. A successful fake-tool CI run is not evidence of video,
audio, input forwarding or compatibility with a specific Android device.

Recorded qualifications:

- [Homatics Dongle G 4K — v0.3.1](qualifications/homatics-dongle-g-4k-v0.3.1.md)
- [Pixel 6 and Xiaomi MiTV-MOOR2 — v0.3.1](qualifications/pixel6-xiaomi-mitv-v0.3.1.md)
