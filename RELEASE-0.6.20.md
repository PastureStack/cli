# CLI 0.6.20

This security maintenance release preserves CLI commands, configuration aliases,
API-key permissions, schema handling and Compose delegation.

- The CLI and maintained source builds use Go 1.27.3 without raising the Go
  language minimum. The official Docker CLI remains in the 29.x series at
  29.9.2, built upstream with patched Go 1.26.10.
- Builder and runtime trust stores use the official October CA update. Installed
  Ubuntu packages receive the reviewed, same-series security fixes; unrelated
  packages and the frozen base snapshot are preserved.
- The four-scope source/product/builder/runtime gate, race tests, reproducible
  builds, SBOMs, exact reviewed vendor policy and immutable asset checks remain
  required. Available fixes, secrets and unregistered findings block release;
  tracking an officially unfixed finding is not a not-affected assertion.

Publication uses the exact successful main-branch producer, validates all archive
checksums, and verifies the published immutable assets anonymously. Numeric
archives retain Server-compatible platform-before-version names, including
`pasturestack-cli-linux-amd64-0.6.20.tar.gz`. Licenses and upstream attribution are
unchanged. Component publication does not establish Server deployment or complete
compatibility-matrix acceptance, and previous tags and assets are not overwritten.
