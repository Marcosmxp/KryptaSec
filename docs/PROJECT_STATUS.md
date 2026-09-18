# Project Status

Last updated: 2026-09-19

## Current phase

**Phase 1 — Core CLI: COMPLETE**

Next planned phase: **Phase 2 — Static application inventory**.

## Phase 1 delivered

- Go 1.27.1 module/toolchain baseline;
- CLI entrypoint and version output;
- configuration loader;
- `kryptasec doctor`;
- runtime/data-directory/SQLite health checks;
- SQLite schema version 1 with future-version rejection;
- local-directory and HTTP target normalization;
- exact-host scope policy with fail-closed behavior;
- cryptographically random scan IDs;
- Phase 1 scan state machine;
- local SQLite scan persistence;
- expected-state transactional scan transitions;
- persisted `kryptasec scan status <id>`;
- JSON structured logging with sensitive-field redaction;
- lifecycle logs that omit full target values;
- unit/integration test baseline;
- GitHub Actions gate for format, tests, vet and CLI build.

## Phase 1 data flow

```text
normalize target
-> INSERT scan(created)
-> transaction: created -> validating_scope
-> exact scope check
-> transaction: validating_scope -> ready | failed
-> later CLI invocation: scan status <id>
```

The system does not yet send HTTP requests, run security scanners, execute commands against a target, invoke an LLM or attempt exploitation.

## Verification

Phase 1 completion code passed GitHub Actions on Go 1.27.1:

```text
gofmt
go test ./...
go vet ./...
go build ./cmd/kryptasec
```

## Project-level work still open

These are governance/release concerns and do not block the Phase 1 technical exit criterion:

- final distribution license decision;
- branch protection/rulesets;
- issue templates;
- automated dependency/license scanning.

## Next engineering milestone

Phase 2 starts with deterministic, read-only local analysis:

1. language/framework detection;
2. dependency inventory;
3. secret-detection adapter;
4. source-tree indexing;
5. configuration discovery;
6. deterministic findings model;
7. JSON/SARIF export.

## Explicitly not started

- LLM agents;
- remote active testing;
- exploitation/validation runtime;
- remediation engine;
- dashboard.

## Current architecture version

`0.1`

Any significant architecture change must update the relevant document and ADR.
