# CLI 0.6.18

This is the standalone CLI, not a Server image or native-acceptance receipt.
It retains the reviewed 0.6.17 schema-URL, JSON flushing and error-propagation
fixes. The older downloaded 0.6.14 client could return exit0 with no output
after a schema403; the maintained inspect action returns that error.

- Signed Ubuntu snapshot `20261015T000000Z`; OpenSSL `3.5.5-1ubuntu3.7`, curl
  `8.18.0-1ubuntu2.8`, OpenSSH client `1:10.2p1-2ubuntu3.8` and applicable base
  fixes replace the older image packages. Go remains `1.27.2`.
- Dapper's official Docker CLI is `29.9.1` (built with Go1.26.9), replacing the
  older Go1.26.5-bearing Docker29.7.2 input. Architecture-specific bytes are pinned.
- Original raw scans remain available. Critical/High, available fixes and
  unregistered findings block; only exact reviewed unfixed Medium/Low records
  use the vendor tracker. No severity-only or blanket build-image exemption.
  The reviewed scan set has 100 exact builder/runtime vendor-pending rows
  across 40 CVEs, with Ubuntu's official status, installed package/PURL,
  upstream CVSS and a 2026-10-20 review date. These are not `not_affected`.
  Separately, 2,691 exact `linux-libc-dev` 7.0.0-38.38 CVE/PURL statements apply
  only to user-space headers: retained dpkg evidence shows no kernel image or
  modules and the shipped Go binary is CGO-disabled. Raw findings are not zero.
- Version, source identity, checksums and reproducible numeric archives must
  pass the release workflow before immutable, non-prerelease publication.
  Publication is not broad CLI or whole-BE18 acceptance.

The executable remains `pasturestack`; configuration aliases, API payloads and
Compose delegation are unchanged. Linux and Darwin amd64 archives and the
Windows amd64 zip retain Server-download-compatible `pasturestack-cli-*` names.
The canonical names place the platform before the version, for example
`pasturestack-cli-linux-amd64-0.6.18.tar.gz`. `SHA256SUMS`, source-producer proof,
raw security reports and SBOMs accompany the checked archives.
The immutable 0.6.14 download and 0.6.17 candidate evidence are not overwritten.
