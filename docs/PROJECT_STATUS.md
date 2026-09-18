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
- local SQLite scan persistence;
- expected-state transactional scan transitions;
- persisted `kryptasec scan status <id>`;
- unit-test baseline.

## Current behavior

```text
normalize target
-> INSERT scan(created)
-> transaction: created -> validating_scope
-> exact scope check
-> transaction: validating_scope -> ready | failed
-> later CLI invocation: scan status <id>
```

The SQLite transition uses the expected previous state and refuses a conflicting update rather than silently overwriting newer state.

The system still does not send HTTP requests, run security scanners, execute commands against a target, or invoke an LLM.

## In progress

- Phase 1 CI verification on Go 1.27.1;
- structured logger wiring;
- database/storage health in `doctor`;
- final distribution license decision;
- branch protection/rulesets.

## Next engineering tasks

1. wire `log/slog` with redaction-safe structured fields;
2. make `doctor` open/check the SQLite database;
3. add schema-version/migration metadata before future schema expansion;
4. complete final Phase 1 build/CI verification;
5. close the Phase 1 exit criterion.

## Explicitly not started

- autonomous pentesting agents;
- remote active testing;
- exploitation/validation runtime;
- remediation engine;
- dashboard.

## Current architecture version

`0.1-draft`

Any significant architecture change must update the relevant document and ADR.
