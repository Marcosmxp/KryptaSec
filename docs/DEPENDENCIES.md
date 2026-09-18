# Dependency Register

Last reviewed: 2026-09-18

## Current production dependencies

| Dependency | Version | License | Purpose |
| --- | --- | --- | --- |
| Go toolchain | 1.27.1 | BSD-3-Clause | compiler, runtime and standard library |
| modernc.org/sqlite | 1.59.0 | BSD-3-Clause | pure-Go SQLite database/sql driver for local scan persistence |

Canonical upstream for `modernc.org/sqlite` is the modernc/cznic SQLite project. The driver is selected because it is pure Go and avoids a CGO dependency for Windows/Linux/macOS CLI distribution.

## GitHub Actions

| Dependency | Pin | Purpose |
| --- | --- | --- |
| actions/checkout | v5 | repository checkout in CI |
| actions/setup-go | v6 | install/pin Go 1.27.1 in CI |

## Policy

Every new direct dependency must update this file in the same pull request and record:

- canonical source;
- pinned version;
- license;
- reason for use;
- relevant security/build implications.
