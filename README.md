# PastureStack CLI

CLI manages compatible environments, hosts, stacks, services, containers, volumes, secrets, catalogs, and accounts through the preserved API.

PastureStack is an independent community effort to preserve, audit, and modernize the Rancher 1.6 ecosystem. It is not affiliated with or endorsed by Rancher Labs or SUSE.

**Upstream:** [`rancher/cli`](https://github.com/rancher/cli). This GitHub fork preserves upstream history, authorship, dates, tags, licenses, and bundled dependency notices. PastureStack maintenance is consolidated into one commit after the preserved upstream boundary.

## Project status

CLI `0.6.18` is maintained from the latest reviewed upstream `v1.6` branch. It uses a digest-pinned Ubuntu 26.04 base, the `20261015T000000Z` signed Ubuntu snapshot, Go 1.27.2, Docker 29.9.1, Go Modules with a checked-in vendor lock, non-root runtime packaging, and reproducible archives. The maintained binary and image name is `pasturestack`.

CLI `0.6.18` retains the reviewed `0.6.17` fix: use the existing schema URL and make `inspect` flush successful JSON and return setup or HTTP errors. This addresses the older downloaded `0.6.14` client's silent success after a denied schema request. Environment URL normalization and Compose delegation remain unchanged. See [release notes](RELEASE-0.6.18.md).

The vendored Logrus path also matches its lowercase import path on case-sensitive Linux filesystems, with file contents and logging behavior unchanged.

The archived `docker/libcompose` source tree and duplicated Compose executor have been removed. The `up` command delegates to the matching maintained `pasturestack-compose` binary; set `PASTURESTACK_COMPOSE_BINARY` only when it is not on `PATH`. API keys are passed through the child environment rather than command-line arguments.

The GitHub security release gate builds and packages twice, runs tests with the race detector, validates source contracts, produces SBOMs, and preserves raw source, binary, build-image and runtime-image scans. Critical/High, available fixes, secrets, and unregistered findings block publication. Only unfixed Medium/Low records may use the exact package/version/CVE tracker with an official vendor source and review date. Reviewed builder kernel-header VEX records do not apply to the shipped runtime. Numeric archives, checksums and source identity are retained for reviewed publication; a green source check alone is not a release or native-acceptance receipt.

## Configuration

Use `PASTURESTACK_URL`, `PASTURESTACK_ACCESS_KEY`, and `PASTURESTACK_SECRET_KEY`. Generic `PLATFORM_*` names, historical environment-variable names, and the existing on-disk configuration are read as migration aliases so established installations and scripts continue to work. Set `PASTURESTACK_LOCALE=en-US` or `zh-TW` for operator messages.

## Build and test

```sh
make build
make test
make package
```

Packaging is local only and does not publish an image. Official numeric release archives use names such as `pasturestack-cli-linux-amd64-0.6.18.tar.gz`; the upstream-boundary `v0.6.14` tag remains immutable. After a normal protected merge and a successful exact-main security release run, `python3 scripts/publish-reviewed-release.py --execute --expected-source <commit> --run-id <run> --receipt <new-file>` publishes the checked bytes as an immutable, non-prerelease `v0.6.18`, with `SHA256SUMS`, raw reports, SBOMs and source-producer proof. Without `--execute` it only prints an offline plan. Publication is not a claim that the whole compatibility or BE18 acceptance matrix has passed. See [COMPATIBILITY.md](COMPATIBILITY.md), [SECURITY.md](SECURITY.md), and [ORIGIN.md](ORIGIN.md).

## License and attribution

The inherited project remains licensed under [Apache License 2.0](LICENSE). Copyright and attribution for inherited work and bundled dependencies remain with their respective authors and contributors. PastureStack contributors claim authorship only for their own changes.
