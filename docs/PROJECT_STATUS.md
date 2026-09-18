# Project Status

Last updated: 2026-09-18

## Current phase

**Phase 0 — Foundation**

The repository is being prepared before implementation of the scanner.

## Completed

- repository created;
- `main` branch established;
- `development` integration branch established;
- CODEOWNERS added;
- basic contribution workflow documented;
- architecture v0.1 drafted;
- security model v0.1 drafted;
- clean-room/provenance policy drafted;
- development workflow drafted;
- roadmap drafted;
- initial ADRs drafted.

## In progress

- project foundation review;
- branch protection configuration;
- final licensing decision.

## Next engineering tasks

1. merge the foundation PR into `development`;
2. select final license/governance model;
3. bootstrap Go module and CLI;
4. add CI for format/test/vet;
5. implement configuration and scan state primitives;
6. implement policy/scope primitives before active network tooling.

## Explicitly not started

- autonomous pentesting agents;
- remote active testing;
- exploitation/validation runtime;
- remediation engine;
- dashboard.

## Current architecture version

`0.1-draft`

Any significant architecture change should update this file and the relevant ADR.
