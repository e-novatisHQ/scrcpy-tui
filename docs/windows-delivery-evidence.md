# Windows delivery evidence and remaining backlog

Snapshot: 2026-09-23. Work branch `codex/windows-v0.5.0`,
[PR #15](https://github.com/e-novatisHQ/scrcpy-tui/pull/15), based on public v0.4.0
(`1ef3ced`). The original plan was uncommitted in the primary checkout and is
preserved by signed commit `e944251`. No v0.5.0 tag or release has been created.
No completed Plane state is inferred from an unmerged implementation.

| Plan item / Plane | Code or deliverable | Verification / limitation | Plane state at snapshot |
|---|---|---|---|
| WIN-01 / TVP-508 | `7843853`, platform session boundary | Linux race, PTY, static checks; Windows cross-build; initial CI 35803402001 | En revue |
| WIN-02 / TVP-509 | `4627d84`, Job Object, plus profile ACL fixes `081d753`, `16a8b22` | Native Windows session and ACL tests passed in job 107002909986; actual scrcpy/Windows 11 pending | En revue |
| WIN-03 / TVP-510 | `590d3f1`, native Go helper fixtures; `2105d74` ACL assertion | Exit, errors, signals, descendants, inherited pipes, outside-process preservation; native race tests passed | En revue |
| WIN-04 / TVP-511 | `1793fbe`, reusable native Windows workflow | Native checks green on bbd8ec8; both exact Windows checks required on main, verified 2026-09-23 | En revue |
| WIN-05 / TVP-512 | `166c9dd`, ZIP, per-target SBOM, checksums, attestation workflow; `3192d7b` normalization | Byte-identical Windows ZIP from Linux and Windows verified; tag-driven draft/attestations not executed | En revue |
| WIN-06 / TVP-513 | `5f5245f`, per-user MSI; `3192d7b`, `c348bc4` verification fixes | Offline payload/tables and native install/upgrade/reinstall/uninstall/PATH cycle passed; Windows 11 GUI pending | En revue |
| WIN-07 / TVC-21 | `b57e285`, [Windows guide](windows.fr.md) | Reviewed against CLI/configuration; final consistency depends on MSI qualification | En validation |
| WIN-08 / LAB-86 | `b57e285`, [Windows 11 protocol](qualifications/windows-11-v0.5.0-protocol.md) | NOT RUN; authorized clean Windows 11 host and Android USB/Wi-Fi access not identified | Bloqué |
| WIN-09 / TVP-514 | [Signing boundary](windows-signing.md) | No configured repository Actions secrets/environments observed; provider/certificate access not supplied | Bloqué |
| CORE-01 / TVP-515 | Capability detection backlog | Depends on v0.5.0; no capability validation implemented in this Windows change | Backlog |
| QUAL-01 / LAB-87 | Supported adb/scrcpy ranges | Depends on CORE-01 and version/device evidence; no new support range claimed | Backlog |
| QUAL-02 / LAB-88 | Debian GUI installation | Clean Debian/Ubuntu desktop authentication and UI recipe still required | Backlog |
| DIST-01 / TVP-516 | RPM assessment | Demonstrated demand and qualified Fedora/RHEL-family runner required | Backlog |
| DIST-02 / TVP-517 | Homebrew assessment | Native macOS support and qualification required before any formula | Backlog |
| REL-01 / TVP-518 | Release-PR automation assessment | Requires additional observed stable manual release cycles | Backlog |
| UX-01 / TVP-519 | Localization assessment | Target locales and scope unarbitrated; current French behavior preserved | Backlog |
| API-01 / TVP-520 | 1.0 contract | Requires several stable multiplatform releases and an agreed observation period | Backlog |

Evidence links: [initial Linux CI](https://github.com/e-novatisHQ/scrcpy-tui/actions/runs/35803402001),
[native Go tests before the ZIP comparison failure](https://github.com/e-novatisHQ/scrcpy-tui/actions/runs/35804683916/job/107002909986),
[current packaging/native pipeline](https://github.com/e-novatisHQ/scrcpy-tui/actions/runs/35805804438).
An overall failed run is not presented as a successful release qualification.

Local checks on the Windows branch: `make check` passed, overall coverage 70.9%
(minimum 70%), `make vuln` found no vulnerabilities, Windows cross-build and vet
passed. The Go skill validator also passed vet/race/build checks for Linux
amd64/arm64. Its generic static contract validator rejects absent `--all`,
`--confirm` and `--json` options; those generic bulk-operation options are not part
of this launcher's accepted CLI, so they were not invented to satisfy the checker.
The project's actual CLI, PTY and session contracts are covered by repository tests.

Reproducibility: Linux archives compare byte-for-byte within the same build
environment; Windows ZIPs also compare byte-for-byte between Linux and Windows
on commit `bbd8ec8`. Per-target SBOMs were generated from a clean local clone
(the SBOM tool cannot resolve this linked worktree). MSI extracted payload and all
installer decision tables compare across two builds; summary UUID/timestamps are
outside that deterministic boundary, as documented in the maintainer guide.

Release blockers remain required review/merge, execution
of the tag draft/attestation path, Windows 11 hardware qualification, and any
signing decision required for publication. Branch protection, signed history,
CODEOWNERS and the existing Linux/security checks must remain enforced.
