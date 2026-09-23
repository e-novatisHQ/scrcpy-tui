# WIN-09 — Authenticode provisioning and release boundary

Status: **blocked on external provisioning; no signing has been performed**.
At the 2026-09-23 inspection, the GitHub repository exposed no configured Actions
secrets or environments. This does not prove that the organization owns no
certificate. The certificate/provider, access method, timestamp service and
responsible owner still need to be identified before a signing implementation can
be selected and tested. Do not put a PFX, password, token or private key in Git,
Plane, a command line, CI artifact or a qualification report.

Required inputs are a suitable code-signing certificate or managed signing
service, an authorized noninteractive authentication method, a trusted RFC 3161
timestamp endpoint, and a clean Windows 11 verification host. The signing account
should be limited to the release identity. Protected CI environment approval must
precede use of that identity; PR code and fork runs must never receive it.

The integration must sign the EXE before producing its ZIP and MSI. Sign the
finished MSI next, then calculate final SHA-256 sums and GitHub attestations.
Do not mutate signed files or publish the unsigned checksum as the signed
checksum. Preserve unsigned reproducibility evidence separately; timestamps and
signature envelopes are not byte-reproducible. The provider-specific signing
command remains intentionally unconfigured until the provisioning inputs exist.

The acceptance record must identify the release commit, public certificate
subject/issuer/thumbprint/expiry, timestamp authority and signed artifact hashes.
On a clean host, verify both EXE and MSI with `Get-AuthenticodeSignature` and the
Windows SDK signature verification tool under the normal trust policy. Require a
valid chain, expected publisher and timestamp; also perform MSI install/upgrade/
uninstall on those exact signed bytes. An unsigned success or a checksum match is
not WIN-09 acceptance. SmartScreen reputation is separate from signature validity.

Rotation: record an owner and expiry monitoring responsibility before enabling
signing; provision a replacement independently, verify test artifacts under both
identities during overlap, then update only the protected signing configuration.
Revoke/disable compromised access at the provider, retain the public provenance
of earlier releases and publish a new version for any replacement artifact.
Never overwrite or retag an existing published release. Validate that logs and
failure artifacts contain no authentication material before enabling the workflow.
