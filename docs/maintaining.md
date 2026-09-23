# Maintainer guide

## Repository settings after creation

Create `e-novatisHQ/scrcpy-tui` as a public repository with `main` as its default
branch. Enable Issues, Discussions if needed, Dependabot alerts/security updates,
secret scanning/push protection where available, and private vulnerability reporting.
Avoid enabling features that are not actively maintained.

Protect `main` with a ruleset:

- Require a pull request, one approval and dismissal of stale approvals.
- Require resolved conversations and the `quality`, `vulnerability`, `linux-amd64`
  and `linux-arm64` checks. Confirm their exact reported names after the first CI run.
- Block force pushes and branch deletion. Limit bypass privileges to maintainers.
- Keep workflow token permissions read-only by default; the release job has its own
  contents-write permission. Do not use `pull_request_target` to execute fork code.

The organization determines who may administer the repository. These settings are
operational steps, not claims that a local file can enforce remote protections.

## Releases

Release hosts need `dpkg-deb` in addition to the development prerequisites.

1. Run `make tools`, `make check`, `make vuln`, `make notices`, and review the
   generated dependency notices.
2. Update `VERSION` and move the curated `Unreleased` entries into a dated section
   in CHANGELOG.md. Commit the reviewed change through a PR.
3. Wait for all main-branch CI checks, including native Linux arm64.
4. Run `make notices sbom release`, verify `dist/SHA256SUMS`, and smoke-test
   extracted archives and Debian packages. Only binaries built on a maintained Go
   toolchain are release candidates.
5. Tag the reviewed commit as `v<contents of VERSION>` and push the tag explicitly.
6. The tag workflow verifies the version/ancestry, reruns quality gates and creates
   a **draft** release. It prepends the curated version section from CHANGELOG.md,
   then GitHub automatically adds categorized merged PRs using `.github/release.yml`.
   Review its assets and generated notes before publishing.

The archives and Debian packages contain the executable, MIT license, README and
generated third-party license notices. Debian packages install the executable under
`/usr/bin`, depend on adb and only suggest scrcpy so that unavailable or obsolete
distro packages do not block installation. Release assets also include a CycloneDX
SBOM. GitHub Actions attests the archives, packages, checksum file and SBOM through
Sigstore; consumers can run
`gh attestation verify <asset> --repo e-novatisHQ/scrcpy-tui`. They do not contain
user configuration, private audit notes,
dev tools, screenshots of real devices, or adb/scrcpy binaries.

Use SemVer. Releases before 1.0 can change public interfaces in a minor version;
record incompatible changes. Never retag a published version. A correction is a new
version. Older binaries can be reinstalled without changing a profile, unless a
future documented migration changes that contract.

Apply `breaking-change`, `enhancement`, `bug`, `security`, `dependencies`, or
`documentation` labels to merged PRs. Unmatched PRs appear under Other changes;
`skip-changelog` excludes maintenance noise. Generated GitHub notes are the detailed
change log, while CHANGELOG.md remains a short curated account of released behavior.

## Dependency and tooling updates

Dependabot covers Go modules and pinned action SHAs. The Makefile separately pins
staticcheck, govulncheck and actionlint: review and bump these explicitly. `.go-version`
is the release/CI toolchain; go.mod is the minimum source language requirement.
Dependabot groups Go updates so coupled libraries are qualified together. Exact
notices are generated in CI and at release time instead of being committed, which
allows dependency PRs to satisfy required checks without stale generated content.
After updates run `go mod tidy`, `make notices`, `make fixtures`, `make check`, and
`make vuln`. Review license changes and module checksum changes before merging.

## Public data hygiene

Use documentation addresses such as `192.0.2.10` and synthetic device IDs. Do not
commit local profiles or real screenshots. `docs/audit-final.md`, `docs/jalons.md`
and `tui.contract.json` are local development records excluded from publication.
Scan the staged tree before the first push; a Git ignore is not a secret scanner.

### Windows gate (v0.5.0 work in progress)

The `windows-amd64` CI job runs native race tests, vet, build and CLI smoke tests
on `windows-2025`. Add this exact check name to protected-main required checks
once it has passed. Preserve `quality`, `vulnerability`, `linux-amd64`,
`linux-arm64`, `codeql`, signed commits, CODEOWNERS review and last-push approval.
A Windows cross-build or hosted runner does not qualify a Windows 11 desktop or
Android hardware. Do not tag/publish v0.5.0 before the delivery plan release gates.

### Windows packaging

`windows-release-check.yml` is shared by PR CI and tag releases. It builds a ZIP
and per-user MSI on Linux, then compares the ZIP with a native Windows rebuild and
runs installation, upgrade, reinstall, PATH preservation and uninstall on the
Windows runner. The synthetic `0.0.1` upgrade fixture is a CI artifact only; it is
never added to release assets. Tag releases include only the native-tested MSI,
verify that its sibling ZIP matches the release ZIP, then checksum and attest it.

The MSI toolchain uses `wixl`/`msitools` 0.106+repack-1 in a digest-pinned Debian
container. `docker build -t scrcpy-tui-msi-toolchain:0.106 packaging/windows`
explicitly provisions it; application startup never downloads tooling. The build
script accepts a ZIP extraction directory containing the executable and notices.
It installs to `%LOCALAPPDATA%\Programs\scrcpy-tui`, uses an HKCU keypath, adds only
that directory to the user PATH, and preserves user profiles on uninstall.
Windows Installer's standard environment actions notify Windows after completion.
Existing terminals still need reopening.

Deterministic boundary: EXE, ZIP and Debian archives are byte-reproducible. The
unsigned MSI payload is the same ZIP executable, and product/component GUIDs are
stable, but msitools writes package UUID and creation/save timestamps at build
time. MSI bytes are therefore **not** claimed reproducible. The final, tested MSI
receives its own SHA-256 and attestation. Authenticode timestamps would introduce
another non-deterministic envelope; verify signed files rather than comparing
signed bytes with an unsigned rebuild. The pinned container base and main tools
are recorded; transitive APT dependencies are resolved at image-build time.

The MSI script overrides wixl 0.106's Environment row with the Windows Installer
`=-PATH` append/remove contract. Its value contains `[~]` so unrelated PATH entries
are preserved. This is checked by native lifecycle tests, not merely XML parsing.
See [the Environment table contract](https://learn.microsoft.com/en-us/windows/win32/msi/environment-table).
