# Changelog

All notable KryptaSec releases will be documented here.

The project has not released a public version yet.

## Unreleased

### Added

- repository governance foundation;
- architecture documentation;
- security model;
- clean-room/provenance policy;
- ADR process;
- initial product roadmap;
- Go 1.27 core CLI bootstrap;
- `kryptasec doctor`;
- configuration via `KRYPTASEC_DATA_DIR` and `KRYPTASEC_LOG_LEVEL`;
- local and HTTP target normalization;
- exact-host authorization scope policy;
- scan IDs and Phase 1 lifecycle state machine;
- SQLite scan persistence;
- SQLite schema version 1 and future-schema rejection;
- SQLite health check in `kryptasec doctor`;
- transactional expected-state lifecycle updates;
- persisted `kryptasec scan status <id>`;
- JSON structured lifecycle logging;
- sensitive structured-field redaction;
- pull-request CI for formatting, tests, vet and build.

### Security

- remote scan preparation fails closed unless the exact hostname is explicitly authorized;
- URL-embedded credentials are rejected;
- subdomains and wildcard patterns are not implicitly authorized;
- concurrent/stale lifecycle transitions fail with a state conflict;
- databases from unsupported newer schema versions fail closed;
- sensitive log fields are redacted;
- lifecycle logs omit full target values;
- Phase 1 scan preparation performs no active network requests.
