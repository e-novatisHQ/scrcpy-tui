# Contributing

Bug reports and focused pull requests are welcome. The interface is currently in
French; documentation and contributions may be in English or French.

## Development

Install Go at the version in `.go-version`, GNU Make, Python 3.11+, and a C compiler
(for the race detector). adb and scrcpy are not required for automated tests.

```sh
make tools
make check
```

Tools are installed under `.tools/`, with versions pinned in the Makefile. This is
an explicit development operation, never performed when starting the launcher.

Keep changes within the existing boundaries: `cmd/` bootstraps the application,
`internal/cli` handles command-line interaction, `internal/app` owns configuration,
ADB discovery, command planning and sessions, and `internal/tui` owns UI state and
rendering. Do not bypass command planning in either frontend.

Before sending a PR:

- Add a regression test for a behavior change, including the failure path.
- Use temporary profiles and fake tools; never test against a contributor's devices.
- Run `make fmt` and `make check`.
- Update documentation and the Unreleased changelog when behavior changes.
- Regenerate fixtures when rendering changes. Dependency notices are generated at
  release time and must not be committed.
- Avoid credentials, real device serials, private addresses, or personal paths.
- Give the PR a concise user-facing title and apply the matching changelog label;
  use `skip-changelog` only for changes with no release-note value.

A PR should describe the trigger, resulting behavior, validation, and limitations.
Changes are reviewed before merging. Dependencies and workflow changes receive
explicit review. No CLA is required; contributions are submitted under MIT.
Participation follows [the code of conduct](CODE_OF_CONDUCT.md).

## Release process

See [the maintainer guide](docs/maintaining.md). Versioning follows SemVer;
pre-1.0 releases may change interfaces in a minor release, documented in the changelog.
GitHub generates categorized release notes from merged PR titles and labels using
`.github/release.yml`; `CHANGELOG.md` keeps the curated, durable release history.
