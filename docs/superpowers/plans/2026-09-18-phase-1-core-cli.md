# Phase 1 Core CLI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Deliver the first safe, executable KryptaSec core: a Go CLI that can validate local targets/scopes, persist scan jobs in SQLite, report environment health, and establish CI without performing active security testing.

**Architecture:** The CLI is a thin adapter over small internal packages. Target normalization and scope authorization are deterministic security boundaries. Scan lifecycle/persistence are independent of future agents so the orchestration layer can be added without changing the CLI contract.

**Tech Stack:** Go 1.27.1, Go standard library (`flag`, `log/slog`), modernc.org/sqlite v1.59.0, GitHub Actions.

**Spec:** `docs/ARCHITECTURE.md`, `docs/ROADMAP.md`, `docs/SECURITY_MODEL.md`

## Global Constraints

- No active network testing in Phase 1.
- Agent/LLM output never makes authorization decisions.
- Local source targets must resolve to an existing directory.
- Remote targets must be HTTP/HTTPS and must pass the explicit scope policy before a scan record is created.
- Secrets must not be persisted in scan records or normal logs.
- Production behavior is implemented test-first.
- `main` is never used for direct development.
- Documentation and project status change in the same PR as implementation.

---

### Task 1: Bootstrap the Go CLI and CI

**Files:**
- Create: `go.mod`
- Create: `cmd/kryptasec/main.go`
- Create: `internal/cli/root.go`
- Create: `internal/cli/root_test.go`
- Create: `.github/workflows/ci.yml`
- Create: `docs/DEPENDENCIES.md`

**Interfaces:**
- Produces: `cli.Run(ctx context.Context, args []string, stdout, stderr io.Writer) int`
- Produces executable entrypoint: `cmd/kryptasec`

- [ ] **Step 1: Write a failing root-command test**

Test that the command is named `kryptasec`, exposes a version, suppresses usage on runtime errors, and contains `doctor` and `scan` once those commands are injected.

- [ ] **Step 2: Run the test and verify RED**

Run:

```bash
go test ./internal/cli -run TestNewRootCommand -v
```

Expected: FAIL because the package/constructor does not exist.

- [ ] **Step 3: Implement the minimal root command**

Use the Go standard library for argument/command routing. Business logic must remain outside `internal/cli`.

- [ ] **Step 4: Add CI**

CI runs on pull requests and pushes to `development`:

```bash
go test ./...
go vet ./...
go build ./cmd/kryptasec
```

- [ ] **Step 5: Verify GREEN and commit**

```bash
go test ./...
go vet ./...
go build ./cmd/kryptasec
```

Commit: `build: bootstrap Go CLI and CI`

---

### Task 2: Configuration and doctor command

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/config_test.go`
- Create: `internal/doctor/check.go`
- Create: `internal/doctor/check_test.go`
- Create: `internal/cli/doctor.go`

**Interfaces:**
- Produces:
  - `config.Load() (config.Config, error)`
  - `doctor.Run(ctx context.Context, cfg config.Config) doctor.Report`

Configuration model:

```go
type Config struct {
    DataDir string
    LogLevel slog.Level
}
```

Doctor report model:

```go
type Check struct {
    Name    string
    OK      bool
    Message string
}

type Report struct {
    Checks []Check
}
```

- [ ] **Step 1: Test config defaults**

Expected data directory is derived from `os.UserConfigDir()` and no API keys are required in Phase 1.

- [ ] **Step 2: Verify RED**

```bash
go test ./internal/config -v
```

- [ ] **Step 3: Implement config**

Support optional environment override:

```text
KRYPTASEC_DATA_DIR
KRYPTASEC_LOG_LEVEL
```

Valid log levels: debug, info, warn, error.

- [ ] **Step 4: Test doctor**

Doctor must report Go runtime metadata and whether the data directory can be created/written. Docker is informational in Phase 1 and must not make the command fail.

- [ ] **Step 5: Implement `kryptasec doctor` and verify**

```bash
go test ./internal/config ./internal/doctor ./internal/cli -v
```

Commit: `feat: add configuration and doctor command`

---

### Task 3: Target normalization and scope policy

**Files:**
- Create: `internal/target/target.go`
- Create: `internal/target/target_test.go`
- Create: `internal/policy/scope.go`
- Create: `internal/policy/scope_test.go`

**Interfaces:**

```go
type Kind string

const (
    KindLocal Kind = "local"
    KindHTTP  Kind = "http"
)

type Target struct {
    Kind      Kind
    Raw       string
    Canonical string
}
```

```go
type Scope struct {
    Hosts []string
}

func (s Scope) Allows(t target.Target) bool
```

- [ ] **Step 1: Write normalization tests**

Required cases:

- existing relative directory -> absolute cleaned local target;
- missing local path -> error;
- `https://EXAMPLE.com:443/a` -> canonical host `example.com`;
- unsupported URL scheme -> error;
- URL credentials -> error.

- [ ] **Step 2: Verify RED**

```bash
go test ./internal/target -v
```

- [ ] **Step 3: Implement target normalization**

Use `net/url`, `filepath.Abs`, `filepath.Clean` and `os.Stat`.

- [ ] **Step 4: Write scope tests**

Required behavior:

- local targets are allowed without remote host scope;
- exact remote host match is allowed;
- sibling/subdomain is denied unless explicitly present;
- host comparison is case-insensitive;
- no wildcard expansion in Phase 1.

- [ ] **Step 5: Implement deterministic scope and verify**

```bash
go test ./internal/target ./internal/policy -v
```

Commit: `feat: add target normalization and scope policy`

---

### Task 4: Scan lifecycle and SQLite persistence

**Files:**
- Create: `internal/scan/scan.go`
- Create: `internal/scan/scan_test.go`
- Create: `internal/store/sqlite/store.go`
- Create: `internal/store/sqlite/store_test.go`

**Interfaces:**

```go
type Status string

const (
    StatusCreated         Status = "created"
    StatusValidatingScope Status = "validating_scope"
    StatusReady           Status = "ready"
    StatusCancelled       Status = "cancelled"
    StatusFailed          Status = "failed"
)

type Scan struct {
    ID        string
    Target    string
    TargetKind string
    Status    Status
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

Store:

```go
type Store interface {
    Create(ctx context.Context, scan scan.Scan) error
    Get(ctx context.Context, id string) (scan.Scan, error)
    UpdateStatus(ctx context.Context, id string, status scan.Status) error
}
```

- [ ] **Step 1: Test legal state transitions**

Phase 1 only needs:

```text
created -> validating_scope -> ready
created/validating_scope/ready -> cancelled
created/validating_scope -> failed
```

Invalid transitions return an error.

- [ ] **Step 2: Verify RED and implement state machine**

```bash
go test ./internal/scan -v
```

- [ ] **Step 3: Write SQLite persistence tests using `t.TempDir()`**

Create/get/update must survive opening a new store against the same DB file.

- [ ] **Step 4: Implement SQLite store**

Use `database/sql` and `modernc.org/sqlite`; schema migration v1 is performed transactionally at open.

- [ ] **Step 5: Verify**

```bash
go test ./internal/scan ./internal/store/sqlite -v
```

Commit: `feat: persist scan lifecycle in SQLite`

---

### Task 5: Implement the safe `scan` command

**Files:**
- Create: `internal/app/service.go`
- Create: `internal/app/service_test.go`
- Create: `internal/cli/scan.go`
- Create: `internal/cli/scan_test.go`

**Interfaces:**

```go
type StartScanRequest struct {
    Target     string
    ScopeHosts []string
}

type StartScanResult struct {
    Scan scan.Scan
}
```

Command contract:

```text
kryptasec scan <target> [--scope-host host]
```

- [ ] **Step 1: Test local scan creation**

A valid local directory creates a persisted scan ending in `ready`.

- [ ] **Step 2: Test remote fail-closed behavior**

A remote URL without matching `--scope-host` returns an error and does not create a ready scan.

- [ ] **Step 3: Implement service orchestration**

Flow:

```text
normalize target
-> create scan(created)
-> validating_scope
-> scope policy
-> ready OR failed
```

No HTTP request is performed.

- [ ] **Step 4: Implement CLI adapter**

Human output must include scan ID, canonical target and status.

- [ ] **Step 5: Verify**

```bash
go test ./...
go vet ./...
go build ./cmd/kryptasec
```

Commit: `feat: add safe scan job creation`

---

### Task 6: Documentation and Phase 1 checkpoint

**Files:**
- Modify: `docs/PROJECT_STATUS.md`
- Modify: `docs/ROADMAP.md`
- Modify: `docs/DEVELOPMENT.md`
- Modify: `CHANGELOG.md`
- Create: `docs/CLI.md`

**Interfaces:**
- Documents actual shipped CLI behavior and remaining Phase 1 work.

- [ ] **Step 1: Document commands**

Include examples for:

```bash
kryptasec doctor
kryptasec scan ./my-app
kryptasec scan https://staging.example.com --scope-host staging.example.com
```

Explicitly state that Phase 1 creates/validates scan jobs but performs no active security testing.

- [ ] **Step 2: Update roadmap checkboxes from implementation evidence**

Do not mark an item complete unless CI/test evidence supports it.

- [ ] **Step 3: Update project status and changelog**

Record branch, architecture version and next tasks.

- [ ] **Step 4: Run final verification**

```bash
go test ./...
go vet ./...
go build ./cmd/kryptasec
```

- [ ] **Step 5: Open PR**

Target: `development`

Title: `feat: bootstrap Phase 1 core CLI`
