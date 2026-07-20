# Security Policy

## Supported state

This repository is under migration review and is not release-ready.

## Security boundaries

- API keys, configuration files, environment variables, registry credentials, Compose secrets, and command output can contain sensitive data.
- Exec, logs, export, and shell commands can expose remote workload data.
- Local configuration must use restrictive permissions and must never be printed with secret values.
- Bundled dependencies and release packaging require review before publication.

## Reporting

Report suspected vulnerabilities through this repository's private security advisory channel. Do not include credentials, private configuration, or production output in a public issue.
