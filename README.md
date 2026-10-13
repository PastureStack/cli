# PastureStack CLI

CLI manages compatible environments, hosts, stacks, services, containers, volumes, secrets, catalogs, and accounts through the preserved API.

PastureStack is an independent community effort to preserve, audit, and modernize the Rancher 1.6 ecosystem. It is not affiliated with or endorsed by Rancher Labs or SUSE.

**Upstream:** [`rancher/cli`](https://github.com/rancher/cli). This GitHub fork preserves upstream history, authorship, dates, tags, licenses, and bundled dependency notices. PastureStack maintenance is consolidated into one commit after the preserved upstream boundary.

## Project status

This is a migration proof of concept based on the latest reviewed upstream `v1.6` branch. Candidate CLI `0.6.17` uses a digest-pinned Ubuntu 26.04 base, the `20260808T000000Z` Ubuntu package snapshot, Go 1.27.0, Docker 29.7.2, Go Modules with a checked-in vendor lock, non-root runtime packaging, and reproducible archives. The maintained binary and image name is `pasturestack`.

CLI `0.6.17` removes a duplicated `/schemas` suffix during client setup and makes `inspect` flush successful JSON and return setup or HTTP errors. Environment URL normalization and Compose delegation remain unchanged. See [release notes](RELEASE-0.6.17.md).

The vendored Logrus path also matches its lowercase import path on case-sensitive Linux filesystems, with file contents and logging behavior unchanged.

The archived `docker/libcompose` source tree and duplicated Compose executor have been removed. The `up` command delegates to the matching maintained `pasturestack-compose` binary; set `PASTURESTACK_COMPOSE_BINARY` only when it is not on `PATH`. API keys are passed through the child environment rather than command-line arguments.

The GitHub security release gate builds and packages twice, runs tests with the race detector, validates source contracts, produces review SBOMs, and scans the source, binary, build image, and runtime image. Successful runs retain the complete numeric candidate archives, checksum manifest, and source revision for reviewed publication. The gate does not create a release, publish an image, or deploy software.

## Configuration

Use `PASTURESTACK_URL`, `PASTURESTACK_ACCESS_KEY`, and `PASTURESTACK_SECRET_KEY`. Generic `PLATFORM_*` names, historical environment-variable names, and the existing on-disk configuration are read as migration aliases so established installations and scripts continue to work. Set `PASTURESTACK_LOCALE=en-US` or `zh-TW` for operator messages.

## Build and test

```sh
make build
make test
make package
```

Packaging is local only and does not publish an image. Candidate CLI archives use the numeric version `0.6.17`; the upstream-boundary `v0.6.14` tag remains immutable. A reviewed numeric release remains a candidate/prerelease, not a stable-quality claim. See [COMPATIBILITY.md](COMPATIBILITY.md), [SECURITY.md](SECURITY.md), and [ORIGIN.md](ORIGIN.md).

## License and attribution

The inherited project remains licensed under [Apache License 2.0](LICENSE). Copyright and attribution for inherited work and bundled dependencies remain with their respective authors and contributors. PastureStack contributors claim authorship only for their own changes.
