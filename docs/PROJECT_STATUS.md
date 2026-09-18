# Project Status

Last updated: 2026-09-18

## Current phase

**Phase 1 — Core CLI (in progress)**

The first executable KryptaSec core is under development on `feature/phase-1-core-cli`.

## Completed

Foundation:
- repository/branch governance;
- CODEOWNERS and contribution workflow;
- architecture v0.1;
- security model v0.1;
- clean-room/provenance policy;
- ADR process;
- CI workflow in `development`.

Phase 1 checkpoint:
- Go module/toolchain baseline;
- CLI entrypoint;
- configuration loader;
- `kryptasec doctor`;
- local-directory and HTTP target normalization;
- exact-host scope policy with fail-closed behavior;
- cryptographically random scan IDs;
- Phase 1 scan state machine;
- `kryptasec scan` validation-only flow;
- unit-test baseline.

## Current behavior

`kryptasec scan` currently performs preparation only:

```text
normalize target
-> create in-memory scan
-> validating_scope
-> exact scope check
-> ready | failed
```

It does not send HTTP requests, run security scanners, execute commands against a target, or invoke an LLM.

## In progress

- Phase 1 CI verification on Go 1.27.1;
- local SQLite persistence;
- structured logger wiring;
- persisted scan retrieval/status commands;
- final distribution license decision;
- branch protection/rulesets.

## Next engineering tasks

1. add the SQLite store behind a small repository interface;
2. persist every scan transition transactionally;
3. add `scan status <id>`;
4. wire `log/slog` with redaction-safe structured fields;
5. make `doctor` report storage/database health;
6. complete the Phase 1 exit criterion and update docs.

## Explicitly not started

- autonomous pentesting agents;
- remote active testing;
- exploitation/validation runtime;
- remediation engine;
- dashboard.

## Current architecture version

`0.1-draft`

Any significant architecture change must update the relevant document and ADR.
