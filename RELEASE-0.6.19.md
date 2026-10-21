# CLI 0.6.19

This compiler-only security maintenance release preserves CLI functionality,
configuration aliases, API-key permissions, schema handling and Compose delegation.

- The CLI compiler is exact Go 1.27.3, including the
  [crypto/x509 fix for CVE-2026-111943](https://pkg.go.dev/vuln/GO-2026-6809).
  Both supported builder architectures retain checksum-pinned official Go archives.
- The builder's official Docker CLI is 29.9.2 with Go 1.26.10, applying the
  [same upstream security fix](https://docs.docker.com/engine/release-notes/29/#2992).
  Architecture-specific archive bytes remain pinned; no vendored behavior is changed.
- The existing four-scope source/product/builder/runtime gate, race tests,
  reproducibility checks, SBOMs and immutable publisher are retained. Critical/High,
  official available fixes, secrets and unregistered records still block release.
  Only exact, currently reviewed, officially unfixed Medium/Low records may be tracked.
  Builder header VEX is not an exemption for the runtime or a Go executable.

Publication requires the complete exact-main release gate and anonymous asset
readback. The numeric archives retain Server-compatible platform-before-version
names, including `pasturestack-cli-linux-amd64-0.6.19.tar.gz`. Original licenses and
upstream attribution remain unchanged. A component release is not Server deployment
or whole compatibility-matrix acceptance; previous tags and assets are not overwritten.
