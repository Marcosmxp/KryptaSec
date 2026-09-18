# KryptaSec Architecture

Status: Draft v0.1  
Last reviewed: 2026-09-19

## 1. Goal

KryptaSec is an autonomous application-security system for authorized assessment of local codebases, test environments and explicitly approved remote targets.

The product should:

1. understand a target and its attack surface;
2. combine deterministic scanners with AI-assisted reasoning;
3. validate findings inside controlled execution environments;
4. preserve evidence and distinguish hypotheses from confirmed vulnerabilities;
5. propose remediations;
6. retest after a remediation;
7. produce machine-readable and human-readable reports.

KryptaSec is an independent implementation. The architecture must not depend on copying source code, prompts, UI, naming, tests or internal implementation details from Strix or other products.

## 2. High-level architecture

```text
                         +----------------------+
                         |   CLI / Local UI     |
                         +----------+-----------+
                                    |
                                    v
                         +----------------------+
                         |   Application Core   |
                         |  Scan Coordinator    |
                         +----------+-----------+
                                    |
              +---------------------+----------------------+
              |                     |                      |
              v                     v                      v
     +----------------+    +----------------+     +----------------+
     | Scope & Policy |    | Target Intake  |     | LLM Gateway    |
     +-------+--------+    +-------+--------+     +-------+--------+
             |                     |                      |
             +----------+----------+-----------+----------+
                        |                      |
                        v                      v
                +---------------+      +---------------+
                | Agent Engine  |      | Tool Registry |
                +-------+-------+      +-------+-------+
                        |                      |
       +----------------+----------------------+----------------+
       |                |                 |                     |
       v                v                 v                     v
 +-----------+    +-----------+    +-----------+        +-----------+
 | Code      |    | Web/API   |    | Browser   |        | Scanners  |
 | Analysis  |    | Testing   |    | Automation|        |/Adapters  |
 +-----+-----+    +-----+-----+    +-----+-----+        +-----+-----+
       |                |                 |                     |
       +----------------+--------+--------+---------------------+
                                 |
                                 v
                         +---------------+
                         | Sandbox Layer |
                         | Docker/Podman |
                         +-------+-------+
                                 |
                                 v
                    +--------------------------+
                    | Evidence & Findings Store|
                    +------------+-------------+
                                 |
                 +---------------+---------------+
                 |                               |
                 v                               v
        +-----------------+             +-----------------+
        | Remediation     |             | Reporting       |
        | / Patch Engine  |             | SARIF/JSON/HTML |
        +-----------------+             +-----------------+
```

## 3. Core components

### 3.1 CLI

Responsibilities:

- parse targets and scan options;
- initialize configuration;
- require/record authorization for active remote testing;
- start, resume, inspect and cancel scans;
- display progress and findings;
- open the local dashboard.

Initial commands:

```text
kryptasec init
kryptasec scan <target>
kryptasec scan status <id>
kryptasec scan cancel <id>
kryptasec findings <scan-id>
kryptasec report <scan-id>
kryptasec config
kryptasec doctor
kryptasec view
```

### 3.2 Application Core / Scan Coordinator

The coordinator owns the scan state machine and must not contain provider-specific or scanner-specific logic.

Proposed states:

```text
created
  -> validating_scope
  -> discovering
  -> analyzing
  -> validating_findings
  -> generating_remediation
  -> retesting
  -> reporting
  -> completed

Any running state -> cancelled | failed
```

### 3.3 Scope & Policy Engine

This is a hard security boundary, not an AI instruction.

Responsibilities:

- canonicalize targets;
- allow or deny domains, IPs, ports and paths;
- prevent scope expansion through redirects or agent decisions;
- apply request budgets and concurrency limits;
- enforce passive/active testing modes;
- log policy decisions.

Agents cannot override policy decisions.

### 3.4 Target Intake

Supported target classes planned for v1:

- local source directory;
- Git repository checked out locally;
- authorized HTTP/HTTPS application;
- OpenAPI specification plus authorized base URL.

Future targets may include containers, mobile packages and cloud configurations.

### 3.5 Agent Engine

Agents are task-specific workers coordinated by the core. The MVP starts with four roles:

- **Scout**: builds an inventory and attack-surface model.
- **Analyst**: correlates code, configuration and runtime observations into test hypotheses.
- **Validator**: attempts controlled validation and produces evidence.
- **Fixer**: proposes minimal patches and requests retesting.

Agent roles are conceptual contracts, not separate permanent processes. The orchestrator may execute them sequentially or concurrently based on policy and resource limits.

### 3.6 LLM Gateway

All model access goes through a provider-neutral interface.

```go
type Provider interface {
    Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
    Stream(ctx context.Context, req CompletionRequest) (<-chan Event, error)
    Capabilities() Capabilities
}
```

Provider adapters should be isolated from the agent engine.

Initial targets:

- OpenAI-compatible APIs;
- Anthropic;
- Google Gemini;
- OpenRouter;
- Ollama/local OpenAI-compatible servers.

Secrets must never be written into scan artifacts.

### 3.7 Tool Registry

Agents invoke capabilities through typed tools rather than arbitrary direct access.

Initial tool classes:

- filesystem read/search;
- source parser;
- package/dependency inventory;
- HTTP client;
- browser driver;
- sandbox shell;
- vulnerability scanners;
- Git diff/patch;
- evidence capture.

Every tool receives a policy context and may be denied before execution.

### 3.8 Sandbox Layer

Potentially unsafe execution occurs outside the main KryptaSec process.

Initial implementation:

- ephemeral Docker container;
- read-only source mount by default;
- explicit writable worktree for patch validation;
- restricted capabilities;
- CPU, memory, process and timeout limits;
- network policy derived from scan scope;
- no Docker socket mounted into worker containers;
- deterministic artifact export.

Podman compatibility is a post-MVP goal.

### 3.9 Evidence and Findings

A finding is not confirmed merely because an LLM predicts it.

Minimum lifecycle:

```text
hypothesis -> candidate -> validated | rejected
```

A validated finding should contain:

- unique ID;
- title and description;
- affected component;
- severity;
- CWE/OWASP mappings when applicable;
- evidence;
- reproducible validation steps;
- remediation;
- validation timestamp;
- tool/agent provenance;
- confidence and validation status.

SQLite is the initial local persistence layer. Export formats must not depend on SQLite internals.

### 3.10 Remediation Engine

The remediation engine:

- creates an isolated writable worktree;
- asks the Fixer for a minimal patch;
- validates formatting/build/tests where available;
- reruns the relevant security check;
- emits a diff;
- never writes to the user's original branch without explicit approval.

Automatic commits and pull requests are future integrations and must remain opt-in.

### 3.11 Reporting

Initial output formats:

- JSON for API/tooling;
- SARIF for code-hosting integrations;
- Markdown for readable local reports;
- HTML for the local dashboard.

PDF is not required for the MVP.

### 3.12 Local Dashboard

React + TypeScript frontend served by the local Go process.

Views:

- scan history;
- scan execution timeline;
- agent/tool activity;
- findings by severity/status;
- evidence;
- remediation diff;
- retest result;
- configuration health.

The dashboard communicates only with the local API in the initial release.

### 3.13 Phase 1 Persistence and Observability

Phase 1 establishes the first implemented orchestration substrate:

- SQLite is the local scan-state store.
- Schema metadata uses `PRAGMA user_version`; the current schema version is `1`.
- Opening a database from a newer unsupported schema fails closed.
- Scan lifecycle updates are transactional and require the expected previous state.
- `kryptasec doctor` validates runtime, data-directory and SQLite health.
- Structured lifecycle logs use Go `log/slog` JSON output.
- Sensitive structured keys are redacted before emission.
- Lifecycle logs identify scans by scan ID and target kind rather than full target path/URL.

Phase 1 intentionally contains no active network testing, scanner execution, LLM invocation or exploitation logic.

## 4. Repository layout

Target layout:

```text
KryptaSec/
├── cmd/
│   └── kryptasec/              # CLI entrypoint
├── internal/
│   ├── app/                    # application services / coordinator
│   ├── agents/                 # Scout/Analyst/Validator/Fixer contracts
│   ├── config/                 # local configuration
│   ├── evidence/               # evidence model and storage
│   ├── findings/               # finding lifecycle and severity
│   ├── llm/                    # provider abstraction + adapters
│   ├── policy/                 # scope/authorization/rate policies
│   ├── remediation/            # patch/retest workflow
│   ├── report/                 # JSON/SARIF/Markdown/HTML
│   ├── runtime/                # sandbox lifecycle
│   ├── target/                 # target normalization/intake
│   └── tools/                  # typed tool registry
├── pkg/
│   └── sdk/                    # stable public Go API, only when needed
├── web/
│   └── dashboard/              # React + TypeScript application
├── containers/                 # sandbox images/policies
├── integrations/
│   └── github/                 # CI/GitHub integration
├── schemas/                    # versioned external schemas
├── docs/
│   ├── adr/
│   └── ...
├── scripts/
├── tests/
│   ├── integration/
│   └── fixtures/
├── .github/
└── go.mod
```

Do not create `pkg/` modules until an API genuinely needs to be public. Default to `internal/`.

## 5. Data flow

### Local source scan

```text
CLI
 -> target normalization
 -> scope/policy validation
 -> workspace snapshot
 -> deterministic discovery
 -> Scout inventory
 -> Analyst hypotheses
 -> Validator controlled tests
 -> evidence store
 -> findings classification
 -> optional Fixer patch
 -> targeted retest
 -> report
```

### Authorized remote scan

```text
CLI
 -> explicit scope declaration
 -> policy compilation
 -> passive discovery
 -> active requests through policy-aware HTTP/browser tools
 -> validation
 -> evidence
 -> report
```

A redirect, discovered hostname or agent suggestion does not automatically expand scope.

## 6. Concurrency model

Go is responsible for orchestration and bounded concurrency.

Rules:

- every scan has a root `context.Context`;
- every tool execution must be cancellable;
- concurrency must be bounded through worker pools/semaphores;
- findings storage must be serialized transactionally;
- external requests must have deadlines;
- cancellation propagates to sandboxes and browser sessions.

## 7. Public interfaces

Version public artifacts independently:

- CLI flags/commands;
- JSON report schema;
- SARIF output;
- configuration file format;
- local API;
- plugin/tool protocol, if introduced.

Breaking changes require an ADR and migration notes once v1.0 is released.

## 8. Non-goals for MVP

- distributed cloud execution;
- autonomous scanning of arbitrary internet assets;
- credential stuffing;
- persistence/post-exploitation;
- malware generation;
- mobile application pentesting;
- Kubernetes-native orchestration;
- enterprise SSO/compliance features.

## 9. Architecture quality gates

A component is ready for implementation only when:

- responsibility and boundaries are documented;
- inputs/outputs are typed;
- threat assumptions are known;
- failure/cancellation behavior is defined;
- tests can be written without real external attack targets;
- third-party dependency provenance is documented.
