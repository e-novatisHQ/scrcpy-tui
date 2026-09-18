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

1. Run `make tools`, `make check`, `make vuln`, and review dependency notices.
2. Update `VERSION` and CHANGELOG.md. Commit the reviewed change through a PR.
3. Wait for all main-branch CI checks, including native Linux arm64.
4. Run `make notices release`, verify `dist/SHA256SUMS`, and smoke-test extracted
   archives. Only binaries built on a maintained Go toolchain are release candidates.
5. Tag the reviewed commit as `v<contents of VERSION>` and push the tag explicitly.
6. The tag workflow verifies the version/ancestry, reruns quality gates and creates
   a **draft** release. Review its assets and release notes before publishing.

The archives contain the executable, MIT license, README and exact third-party
license notices. They do not contain user configuration, private audit notes,
dev tools, screenshots of real devices, or adb/scrcpy binaries.

Use SemVer. Releases before 1.0 can change public interfaces in a minor version;
record incompatible changes. Never retag a published version. A correction is a new
version. Older binaries can be reinstalled without changing a profile, unless a
future documented migration changes that contract.

## Dependency and tooling updates

Dependabot covers Go modules and pinned action SHAs. The Makefile separately pins
staticcheck, govulncheck and actionlint: review and bump these explicitly. `.go-version`
is the release/CI toolchain; go.mod is the minimum source language requirement.
After updates run `go mod tidy`, `make notices`, `make fixtures`, `make check`, and
`make vuln`. Review license changes and module checksum changes before merging.

## Public data hygiene

Use documentation addresses such as `192.0.2.10` and synthetic device IDs. Do not
commit local profiles or real screenshots. `docs/audit-final.md`, `docs/jalons.md`
and `tui.contract.json` are local development records excluded from publication.
Scan the staged tree before the first push; a Git ignore is not a secret scanner.
