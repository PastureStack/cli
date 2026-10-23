# PastureStack CLI

CLI manages compatible environments, hosts, stacks, services, containers, volumes, secrets, catalogs, and accounts through the preserved API.

PastureStack is an independent community effort to preserve, audit, and modernize the Rancher 1.6 ecosystem. It is not affiliated with or endorsed by Rancher Labs or SUSE.

**Upstream:** [`rancher/cli`](https://github.com/rancher/cli). This GitHub fork preserves upstream history, authorship, dates, tags, licenses, and bundled dependency notices. PastureStack maintenance is consolidated into one commit after the preserved upstream boundary.

## Project status

Published CLI [`0.6.20`](https://github.com/PastureStack/cli/releases/tag/v0.6.20) is maintained from the latest reviewed upstream `v1.6` branch. It uses Go 1.27.3, the checksum-pinned official Docker 29.9.2 builder input, Go Modules with a checked-in vendor lock, and reproducible archives. Runtime packaging uses a digest-pinned Ubuntu 26.04 base, the `20261015T000000Z` signed Ubuntu snapshot and a non-root user. It is CI-qualified, not a published registry image. The maintained command name is `pasturestack`.

The immutable, non-prerelease archives were produced from signed, normally merged source [`60a13c6ef8e60018e6409fb7b510bcf0d44c5030`](https://github.com/PastureStack/cli/commit/60a13c6ef8e60018e6409fb7b510bcf0d44c5030), after the complete [exact-main release gate](https://github.com/PastureStack/cli/actions/runs/39679384298). All 26 assets were anonymously downloaded and matched their recorded bytes and SHA-256. Later documentation commits are not a different producer and do not rebuild or overwrite these assets.

CLI `0.6.20` applies the Go crypto/x509 security fix without changing CLI functionality or API-key permissions. `inspect` uses the existing schema URL, flushes successful JSON and returns setup or HTTP errors. Environment URL normalization and Compose delegation remain unchanged. See [release notes](RELEASE-0.6.20.md).

The vendored Logrus path also matches its lowercase import path on case-sensitive Linux filesystems, with file contents and logging behavior unchanged.

The archived `docker/libcompose` source tree and duplicated Compose executor have been removed. The `up` command delegates to the matching maintained `pasturestack-compose` binary; set `PASTURESTACK_COMPOSE_BINARY` only when it is not on `PATH`. API keys are passed through the child environment rather than command-line arguments.

The GitHub security release gate builds and packages twice, runs tests with the race detector, validates source contracts, produces SBOMs, and preserves raw source, binary, build-image and runtime-image scans. Critical/High, available fixes, secrets, and unregistered findings block publication. Only unfixed Medium/Low records may use the exact package/version/CVE tracker with an official vendor source and review date. Reviewed builder kernel-header VEX records do not apply to the shipped runtime. Numeric archives, checksums and source identity are retained for reviewed publication; a green source check alone is not a release or native-acceptance receipt.

For this release, source and shipped-binary vulnerability scans are zero; raw builder/runtime findings are not zero. Exact findings, current vendor-review dates and upstream CVSS are retained in the immutable release's [tracker](https://github.com/PastureStack/cli/releases/download/v0.6.20/cli-vendor-pending.json) and [security summary](https://github.com/PastureStack/cli/releases/download/v0.6.20/cli-security-summary.json). Builder `linux-libc-dev` header VEX rows are supported by retained package inventory (no kernel image/modules) and CGO-disabled binary metadata; they do not exempt runtime packages or executable code. Native CLI/Server interoperability and the broader compatibility matrix remain separately evidenced; publication alone is not a whole-matrix PASS.

## Configuration

Use `PASTURESTACK_URL`, `PASTURESTACK_ACCESS_KEY`, and `PASTURESTACK_SECRET_KEY`. Generic `PLATFORM_*` names, historical environment-variable names, and the existing on-disk configuration are read as migration aliases so established installations and scripts continue to work. Set `PASTURESTACK_LOCALE=en-US` or `zh-TW` for operator messages.

## Build and test

```sh
make build
make test
make package
```

Packaging is local only and does not publish an image. Official numeric release archives use names such as `pasturestack-cli-linux-amd64-0.6.20.tar.gz`; the upstream-boundary `v0.6.14` tag remains immutable. The reviewed `scripts/publish-reviewed-release.py` entry verifies the normal protected source merge, required checks and exact-main gate; it publishes source-bound checked bytes with `SHA256SUMS`, raw reports, SBOMs and producer proof, then anonymously verifies them. Without `--execute` it only prints an offline plan. An existing tag is never overwritten. See [COMPATIBILITY.md](COMPATIBILITY.md), [SECURITY.md](SECURITY.md), and [ORIGIN.md](ORIGIN.md).

## License and attribution

The inherited project remains licensed under [Apache License 2.0](LICENSE). Copyright and attribution for inherited work and bundled dependencies remain with their respective authors and contributors. PastureStack contributors claim authorship only for their own changes.

## October 2026 security maintenance

Both builder and runtime apply Ubuntu's fixed `zlib1g 1:1.3.dfsg+really1.3.1-1ubuntu3.2` package for CVE-2026-85091. The fixed finding is removed from the vendor-pending register rather than suppressed.

Maintained Go builds select the patched Go 1.27.3 compiler without raising the existing Go language requirement. Ubuntu build/runtime stages apply the checksum-pinned `ca-certificates 20260816ubuntu0.26.04.1` security package after their frozen APT baseline; the final trust store and installed version are verified. Release artifacts must pass the repository's existing checks before publication.
