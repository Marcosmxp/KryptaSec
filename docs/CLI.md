# KryptaSec CLI

Status: Phase 1 checkpoint  
The CLI contract documented here reflects implemented behavior, not future roadmap features.

## Build

KryptaSec uses Go 1.27.1 as the pinned toolchain.

```bash
go build -o kryptasec ./cmd/kryptasec
```

On Windows:

```powershell
go build -o kryptasec.exe ./cmd/kryptasec
```

## Version

```bash
kryptasec --version
```

Development builds report `dev` unless the version variable is set at build time.

## Doctor

```bash
kryptasec doctor
```

Current checks:

- Go runtime/OS/architecture;
- KryptaSec data-directory creation and write access.

Environment:

```text
KRYPTASEC_DATA_DIR
KRYPTASEC_LOG_LEVEL=debug|info|warn|error
```

## Scan preparation

### Local directory

```bash
kryptasec scan ./my-app
```

Requirements:

- the path must exist;
- the path must be a directory;
- it is canonicalized to an absolute clean path.

### Remote HTTP/HTTPS target

```bash
kryptasec scan --scope-host staging.example.com https://staging.example.com
```

The exact hostname must be provided through `--scope-host`.

Phase 1 deliberately does not infer authorization:

- `example.com` does not authorize `api.example.com`;
- `*.example.com` is not expanded;
- redirects do not expand scope;
- credentials embedded in URLs are rejected.

Multiple exact hosts can be provided by repeating the flag:

```bash
kryptasec scan \
  --scope-host app.example.com \
  --scope-host api.example.com \
  https://app.example.com
```

## Current scan lifecycle

```text
created
  -> validating_scope
  -> ready

created/validating_scope
  -> failed

created/validating_scope/ready
  -> cancelled
```

The current command stops at `ready`.

## Security boundary

At this checkpoint, `kryptasec scan`:

- does not send an HTTP request;
- does not crawl the target;
- does not execute shell commands;
- does not invoke an LLM;
- does not exploit vulnerabilities;
- does not persist the scan yet.

SQLite persistence and persisted status commands are the next Phase 1 implementation step.
