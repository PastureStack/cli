# PastureStack CLI

CLI manages compatible environments, hosts, stacks, services, containers, volumes, secrets, catalogs, and accounts through the preserved API.

PastureStack is an independent community effort to preserve, audit, and modernize the Rancher 1.6 ecosystem. It is not affiliated with or endorsed by Rancher Labs or SUSE.

**Upstream:** [`rancher/cli`](https://github.com/rancher/cli). This GitHub fork preserves upstream history, authorship, dates, tags, licenses, and bundled dependency notices. PastureStack maintenance is consolidated into one commit after the preserved upstream boundary.

## Project status

This is a migration proof of concept based on the latest reviewed upstream `v1.6` branch. Candidate `0.6.15` uses a digest-pinned Ubuntu 26.04 base, the `20260808T000000Z` Ubuntu package snapshot, Go 1.26.5, Docker 29.6.2, non-root runtime packaging, and reproducible archives. The maintained binary and image name is `pasturestack`.

The archived `docker/libcompose` source tree has been removed. The small conversion, logging, and Compose-label contracts still required by the bundled executor are maintained and tested locally in this repository.

The GitHub security release gate builds and packages twice, runs tests with the race detector, validates source contracts, produces review SBOMs, and scans the source, binary, build image, and runtime image. It does not publish binaries or images, create a release, or deploy software.

## Configuration

Use `PASTURESTACK_URL`, `PASTURESTACK_ACCESS_KEY`, and `PASTURESTACK_SECRET_KEY`. Generic `PLATFORM_*` names, historical environment-variable names, and the existing on-disk configuration are read as migration aliases so established installations and scripts continue to work. Set `PASTURESTACK_LOCALE=en-US` or `zh-TW` for operator messages.

## Build and test

```sh
make build
make test
make package
```

Packaging is local only and does not publish an image. Candidate archives use the numeric version `0.6.15`; the upstream-boundary `v0.6.14` tag remains immutable. See [COMPATIBILITY.md](COMPATIBILITY.md), [SECURITY.md](SECURITY.md), and [ORIGIN.md](ORIGIN.md).

## License and attribution

The inherited project remains licensed under [Apache License 2.0](LICENSE). Copyright and attribution for inherited work and bundled dependencies remain with their respective authors and contributors. PastureStack contributors claim authorship only for their own changes.
