# Compatibility Contract

CLI preserves API resource and action names, Compose and catalog fields, established environment-variable aliases, on-disk configuration migration, generated client types, server-provided links, and remote error data.

The product-facing executable, image, source path, help text, and new settings use PastureStack naming. Historical identifiers remain only as read-only or fallback compatibility contracts for existing scripts and servers.

The build now uses Go Modules and a checked-in vendor lock. The generated Rancher 1.6 API schema bindings are protocol compatibility code, not a maintained upstream dependency; the exact subset required by this product is retained under `internal/rancherclient` with its Apache-2.0 provenance. Its HTTP and WebSocket transports use maintained dependencies.

The archived `docker/libcompose` tree and the duplicated legacy executor are no longer part of this repository. Focused conversion, logger, and Compose-label contracts remain under `composecompat`. Full Compose reconciliation belongs to the dedicated `pasturestack-compose` product; `pasturestack up` delegates to that binary and preserves established flags and environment aliases.

Operator lifecycle messages support `en-US` and `zh-TW`. Resource names, API payloads, Compose content, identifiers, and remote errors are not translated.

Before release, validate login, context and config migration, list and inspect, create and update, delete, execute, logs, Compose import/export, catalog, secrets, volumes, shell completion, exit codes, and rollback against an isolated compatible server.
