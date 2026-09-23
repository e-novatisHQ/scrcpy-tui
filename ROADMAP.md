# Roadmap

Delivered capabilities and atomic work items are tracked in the
[project status and delivery plan](docs/project-status-and-plan.md).

## Near term

- Deliver a native Windows 10/11 x64 build with owned-process cleanup, native CI,
  a verified ZIP, a click-to-install MSI and real-host qualification.
- Qualify releases against a documented real-device matrix.
- Define supported adb and scrcpy version ranges from recorded evidence.
- Add regression coverage for defects reported by users.
- Evaluate a release-PR bot after several manual SemVer cycles.

## Later

- Assess additional packaging based on demonstrated demand (RPM or Homebrew after
  their host platforms are supported).
- Assess Windows ARM64 after adb, scrcpy and hardware availability are demonstrated.
- Evaluate interface localization without weakening the current French UX.
- Consider 1.0 after the configuration and CLI contracts remain stable across
  several releases.

Roadmap items are intentions, not delivery commitments. Track accepted work in
GitHub Issues and link it to reproducible evidence.
