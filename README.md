# scrcpy-tui

A fast terminal launcher for [scrcpy](https://github.com/Genymobile/scrcpy): choose
an ADB device, choose a preset, inspect the command and press Enter.

The interface is currently in French. [Guide en français](docs/usage.fr.md).

## Features

- ADB inventory with model information; offline/unauthorized devices are disabled.
- Device search, paging, persistent selections and a responsive keyboard interface.
- Three editable presets and free scrcpy arguments, without invoking a shell.
- Exact command preview, pre-launch device revalidation and session diagnostics.
- Atomic local preset storage and recovery copies for invalid configurations.
- Foreground sessions: the menu returns when scrcpy closes.
- Noninteractive inventory, preset listing, preview and explicit launch commands.

## Requirements and support

**Linux** with adb and scrcpy installed in `PATH`. The build does not bundle these
external tools or install anything when launched. For Wi-Fi, connect your device
through adb first. Native amd64 is tested locally; native arm64 is configured in
CI and remains pending until the public CI has run. macOS and Windows are not
supported in this release. No graphical/video qualification is claimed by fake-tool tests.

## Install

After a release is published, download the Linux archive matching your CPU from
[GitHub Releases](https://github.com/e-novatisHQ/scrcpy-tui/releases), together with
`SHA256SUMS`. Verify the archive before extracting it:

```sh
# Example for VERSION 0.3.0; choose linux_arm64 for arm64.
sha256sum --ignore-missing --check SHA256SUMS
mkdir scrcpy-tui-release
tar -xzf scrcpy-tui_0.3.0_linux_amd64.tar.gz -C scrcpy-tui-release
install -Dm755 scrcpy-tui-release/scrcpy-tui "$HOME/.local/bin/scrcpy-tui"
scrcpy-tui
```

No binary release has been published yet. To build from a local checkout:

```sh
make build
./bin/scrcpy-tui
```

Use Go at the version in `.go-version` for development and releases. This builds a
single executable; Go is not required on the target machine. `make install` installs
under `~/.local/bin`; set `PREFIX` to choose another prefix.

## Keyboard

| Key | Action |
|---|---|
| Tab / ← / → | Switch between device and preset lists |
| ↑ / ↓ | Select; disabled devices are skipped |
| Home / End, Page↑ / Page↓ | Jump or page without wrapping |
| / | Filter devices by serial, address, port or model |
| Enter | Finish search, or launch the selected device/preset |
| Escape | Clear the filter or close/cancel a modal |
| n / e / d | Create / edit / delete the selected preset, regardless of focus |
| c / l / ? | Full command / session diagnostics / keyboard help |
| r / q / Ctrl+C | Refresh / quit / quit the menu |

The preset editor uses Tab/Shift+Tab for fields, Enter to save, Escape to cancel.
Deletion requires `o`; any other key cancels. Ctrl+C during scrcpy stops that session
and returns to the menu. SIGTERM stops the owned session and exits the launcher.

## CLI

Options work before or after the subcommand:

```sh
scrcpy-tui devices
scrcpy-tui presets
scrcpy-tui preview --device USB123 --preset 'Très léger'
scrcpy-tui launch --device USB123 --preset 'Léger Wi-Fi' --yes
scrcpy-tui --help
```

`--args "--window-title 'My TV'"` adds arguments. Unknown scrcpy options remain
allowed; device selectors are reserved for the launcher. Preview is read-only.
Exit codes: 0 success/menu quit, 1 technical error, 3 invalid input, 4 missing CLI
`--yes`, 130 interrupted CLI session.

## Configuration

`$XDG_CONFIG_HOME/scrcpy-tui/presets.json`, or `~/.config/scrcpy-tui/presets.json`.
Use `--config /path/profile.json` for an isolated profile. Reads do not write.
Files and recovery copies use mode 0600. Before replacing an invalid configuration,
original bytes are retained in `presets.json.recovery-*`. Close the launcher before
restoring a backup or editing the JSON manually. Concurrent profile writes are not merged.

Diagnostics retain the last 32 KiB of scrcpy output in memory, with no persistent
log. Arguments appear in previews, profiles and process listings: do not put secrets
in presets. Options passed to scrcpy keep their own effects.

## Development

```sh
make tools          # Explicitly install pinned tools under .tools/
make check          # Format, vet, race tests, staticcheck, actionlint, PTY tests
make vuln           # govulncheck; needs access to the public vulnerability DB
make notices        # Regenerate exact dependency license notices
make fixtures       # Synthetic text renderings
make release        # Linux amd64/arm64 archives and SHA256SUMS
```

GNU Make, Python 3.11+ and a C compiler are needed for the complete validation.
Tests use temporary profiles and fake devices. See [CONTRIBUTING](CONTRIBUTING.md),
[architecture](docs/architecture.md), [maintaining](docs/maintaining.md), and
[security](SECURITY.md). Releases are tagged, quality-gated and created as drafts.

## License

[MIT](LICENSE), © e-novatisHQ and contributors. Dependencies retain their own
licenses; see [third-party notices](THIRD_PARTY_NOTICES.md). scrcpy-tui is an
independent project, not affiliated with Genymobile or Google.
