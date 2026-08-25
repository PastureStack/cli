# Internal Docker API proxy

This package preserves the Rancher v1 Docker socket tunnelling protocol used by
the PastureStack CLI. It is internal implementation code, not a reusable public
module.

The source was imported from `github.com/rancher/rancher-docker-api-proxy` at
the repository's pinned legacy revision and is maintained here so the CLI no
longer depends on an abandoned GOPATH package graph. Its original Apache-2.0
license remains in this repository.
