# Changelog

## Unreleased

## 0.3.1 — 2026-09-20

- Add CodeQL analysis, a coverage floor, grouped dependency updates and CODEOWNERS.
- Generate CycloneDX SBOMs and signed provenance attestations for release assets.
- Verify release reproducibility in CI and generate dependency notices at release time.
- Add support, conduct, roadmap and real-device qualification documentation.

## 0.3.0 — 2026-09-20

- Update Bubbles to 1.0.0, x/ansi to 0.11.8, go-shellwords to 1.0.15 and
  x/term to 0.46.0; the minimum Go version is now 1.26.
- Public Go module and testable CLI shared with a small executable entrypoint.
- Automated formatting, race tests, vet, static analysis and vulnerability checks.
- Isolated PTY integration tests, Linux release archives and checksum verification.
- MIT licensing, dependency notices and contributor/maintainer documentation.
- CI and release workflows with pinned actions and limited permissions.

## 0.2.1

- Preserve invalid configurations in a recovery copy before replacement.
- Relay session signals and clean up owned process groups and inherited output.
- Prevent concurrent discovery and launch preparation.
- Clamp page navigation and keep errors visible in small terminals.
- Reject alternate device selectors and respect the CLI `--` separator.

## 0.2.0

- Device filtering, paging, focus styling, responsive panels and command details.
- Built-in keyboard help and bounded scrcpy diagnostics.
- Options accepted before or after CLI subcommands.

## 0.1.0

- Initial ADB inventory, preset editor, persistence, preview and foreground launch.
