# CLI 0.6.17 candidate

This version belongs to PastureStack CLI, not the Server runtime. It is not a
publication or native-acceptance receipt. The immutable upstream `v0.6.14` tag
is unchanged. Publication, if approved, remains a candidate/prerelease; the
broader CLI compatibility contract and migration-review state are unchanged.

- Use the environment's existing schema URL directly; do not append a second
  `/schemas` segment.
- Flush successful `inspect` output and return setup, HTTP, or output errors.
- Preserve `EnvironmentURL`, `NormalizeUrl`, API payloads, server-provided links,
  configuration aliases, and `up` delegation to the maintained Compose binary.

Focused validation selector:

```sh
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -mod=vendor ./cmd -run '^TestCLISchemaInspect017' -count=1
```

The five focused tests use the real client and native inspect action against
local `httptest` responses. They cover the schema request URL, JSON flushing,
missing configuration, schema HTTP denial, and resource HTTP denial. They do
not constitute live Server, CLI publication, or BE-18 whole-row acceptance.

The release gate packages numeric CLI version `0.6.17`. For a local candidate
build, use `VERSION_OVERRIDE=0.6.17`; source scripts otherwise retain their
existing tag-derived version behavior.

After a successful complete gate, the artifact
`cli-release-candidate-0.6.17-<source revision>` retains all files under
`dist/artifacts/0.6.17/`, `evidence/artifacts.sha256`, and the source revision
for 14 days. Reviewed publication must use these actual checked bytes, verify
the checksum manifest and tag/source binding, and read back the uploaded
release assets. The gate itself does not publish a release or bypass required
repository checks.
