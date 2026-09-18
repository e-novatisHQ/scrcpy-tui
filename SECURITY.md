# Security policy

Security fixes are provided for the latest release. Linux is the supported platform.

Please report suspected vulnerabilities privately using the repository's
**Security → Report a vulnerability** feature. Do not open a public issue containing
credentials, device identifiers, logs from a private network, or an exploit.
If private reporting is unavailable, contact the organization privately through
its GitHub profile before sharing details. No response-time SLA is currently offered.

The launcher does not invoke a shell for preset arguments. Device selection is
validated again before launch. Only its own session process group is signalled.
Presets remain arbitrary scrcpy arguments and retain scrcpy's own effects.
Local configuration and recovery copies use restrictive permissions. Command-line
arguments and previews are visible to users of the machine; do not store secrets
in presets. adb and scrcpy are separately installed, trusted external dependencies.

GitHub workflows use immutable action references, read-only default permissions,
and no privileged execution of fork pull requests. A release job has write access
only to publish explicitly tagged releases.
