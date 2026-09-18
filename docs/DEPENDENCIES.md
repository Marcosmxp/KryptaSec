# Dependency Register

Last reviewed: 2026-09-18

## Current production dependencies

The current Phase 1 checkpoint uses only the Go standard library at runtime.

| Dependency | Version | License | Purpose |
| --- | --- | --- | --- |
| Go toolchain | 1.27.1 | BSD-3-Clause | compiler, runtime and standard library |

## GitHub Actions

| Dependency | Pin | Purpose |
| --- | --- | --- |
| actions/checkout | v5 | repository checkout in CI |
| actions/setup-go | v6 | install/pin Go 1.27.1 in CI |

## Planned, not yet imported

SQLite persistence is planned for the next Phase 1 increment. The selected driver must be pure Go, actively maintained and license-compatible. It must not be added until its version/license/provenance is recorded here and its tests are present.

## Policy

Every new direct dependency must update this file in the same pull request and record:

- canonical source;
- pinned version;
- license;
- reason for use;
- relevant security/build implications.
