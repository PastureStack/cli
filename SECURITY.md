# Security Policy

## Supported state

This repository is under migration review and is not release-ready.

## Security boundaries

- API keys, configuration files, environment variables, registry credentials, Compose secrets, and command output can contain sensitive data.
- Exec, logs, export, and shell commands can expose remote workload data.
- Local configuration must use restrictive permissions and must never be printed with secret values.
- Bundled dependencies and release packaging require review before publication.
- Candidate changes must pass the GitHub security release gate for source, binary, build-image, runtime-image, secret, reproducibility, and SBOM evidence.

## Release CVE policy

Critical/High findings, findings with an available fixed version, secrets, and
unregistered records block release. Only Medium/Low findings without an official
fixed version may be tracked: exact component, CVE, PURL, installed version,
official vendor status/source, verification time and a non-expired review date
are required. A tracker is not a `not_affected` assertion. Raw scan reports are
retained. Builder-only kernel-header VEX records require no installed kernel
image/modules and must match the observed CVE/PURL set; they never exempt the
runtime or a vulnerable tool binary. Applied Ubuntu curl fixes use `fixed` VEX
with the official backport evidence, not a fabricated `not_affected` status.

## Reporting

Report suspected vulnerabilities through this repository's private security advisory channel. Do not include credentials, private configuration, or production output in a public issue.
