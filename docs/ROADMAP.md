# KryptaSec Roadmap

Status: Draft v0.1  
Roadmap dates are intentionally milestone-based rather than calendar promises.

## Phase 0 — Foundation

Goal: make the repository safe and maintainable before product code.

- [x] `main` and `development` branch model
- [x] CODEOWNERS
- [x] contribution policy
- [x] architecture v0.1
- [x] clean-room/provenance policy
- [x] security model
- [x] ADR process
- [ ] choose final distribution license
- [ ] branch protection/rulesets
- [ ] issue/PR templates
- [x] CI foundation
- [ ] dependency/license scanning

Exit criterion: repository governance and architecture are agreed and enforced.

## Phase 1 — Core CLI — COMPLETE

Goal: a reliable local executable with no autonomous offensive behavior yet.

- [x] Go module and CLI bootstrap
- [x] configuration subsystem
- [x] structured logging
- [x] scan IDs and lifecycle
- [x] local SQLite storage
- [x] target normalization
- [x] scope policy parser
- [x] `kryptasec doctor`
- [x] database health check
- [x] schema versioning/migration metadata
- [x] unit/integration test baseline
- [x] persisted `scan status <id>`
- [x] CI format/test/vet/build gate

Exit criterion met: `kryptasec scan` creates, transactionally tracks and retrieves a safe analysis job while active network testing remains disabled.

Completed: 2026-09-19.

## Phase 2 — Static application inventory

- [ ] language/framework detection
- [ ] dependency inventory
- [ ] secret-detection adapter
- [ ] source tree indexing
- [ ] configuration discovery
- [ ] deterministic findings API
- [ ] JSON/SARIF output

Exit criterion: useful non-LLM local security inventory.

## Phase 3 — LLM and agent foundation

- [ ] provider-neutral LLM interface
- [ ] OpenAI-compatible adapter
- [ ] Ollama/local adapter
- [ ] structured agent messages
- [ ] typed tool registry
- [ ] Scout agent
- [ ] Analyst agent
- [ ] token/cost accounting
- [ ] redaction of secrets before model calls

Exit criterion: agents can reason about a local codebase using policy-controlled read-only tools.

## Phase 4 — Sandbox and validation

- [ ] Docker sandbox
- [ ] resource/time limits
- [ ] scope-aware network policy
- [ ] HTTP testing tool
- [ ] Playwright browser tool
- [ ] Validator agent
- [ ] evidence capture
- [ ] finding lifecycle: hypothesis/candidate/validated/rejected

Exit criterion: the system can safely validate selected vulnerability classes in owned/test fixtures.

## Phase 5 — Remediation

- [ ] isolated writable worktree
- [ ] Fixer agent
- [ ] patch generation
- [ ] build/test hooks
- [ ] targeted security retest
- [ ] human approval gate
- [ ] patch export

Exit criterion: confirmed fixture vulnerability can be patched and revalidated end-to-end.

## Phase 6 — Dashboard

- [ ] React + TypeScript workspace
- [ ] local Go API
- [ ] scan history
- [ ] live event stream
- [ ] findings/evidence UI
- [ ] patch/retest UI
- [ ] accessibility baseline

Exit criterion: complete local workflow can be operated without reading raw artifacts.

## Phase 7 — CI/CD and integrations

- [ ] headless mode
- [ ] deterministic exit codes
- [ ] GitHub Actions example/action
- [ ] SARIF upload workflow
- [ ] diff-aware scanning
- [ ] optional PR remediation flow

Exit criterion: KryptaSec can gate a test repository safely in CI.

## Phase 8 — v1.0 hardening

- [ ] compatibility matrix
- [ ] threat-model review
- [ ] reproducible releases
- [ ] signed checksums/releases
- [ ] migration/versioning policy
- [ ] documentation audit
- [ ] benchmark fixtures
- [ ] security disclosure process
- [ ] external dependency/license audit

Exit criterion: stable local-first v1.0.
