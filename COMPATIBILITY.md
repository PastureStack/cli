# Compatibility Contract

CLI preserves API resource and action names, Compose and catalog fields, established environment-variable aliases, on-disk configuration migration, generated client types, server-provided links, and remote error data.

The product-facing executable, image, source path, help text, and new settings use PastureStack naming. Historical identifiers remain only as read-only or fallback compatibility contracts for existing scripts and servers.

The checked-in Go vendor tree retains upstream module import paths, package declarations, source attribution, copyright notices, and licenses. Those paths identify third-party source and are required by the pre-module Go build graph; rewriting them would misstate provenance and break imports. They are not PastureStack product branding.

Operator lifecycle messages support `en-US` and `zh-TW`. Resource names, API payloads, Compose content, identifiers, and remote errors are not translated.

Before release, validate login, context and config migration, list and inspect, create and update, delete, execute, logs, Compose import/export, catalog, secrets, volumes, shell completion, exit codes, and rollback against an isolated compatible server.
