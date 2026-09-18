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

Multiple exact hosts can be provided by repeating the flag.

## Persisted status

Each accepted or rejected scan job is stored in:

```text
<KRYPTASEC_DATA_DIR>/kryptasec.db
```

When `KRYPTASEC_DATA_DIR` is unset, the OS user configuration directory is used.

Read a scan later with:

```bash
kryptasec scan status <scan-id>
```

The command returns the canonical target, target kind, current status, creation time and last update time.

## Transactional lifecycle

```text
created
  -> validating_scope
  -> ready

created/validating_scope
  -> failed

created/validating_scope/ready
  -> cancelled
```

Each transition is persisted in a SQLite transaction and includes the expected previous status. If the stored state no longer matches the expected state, the transition fails with a conflict rather than overwriting newer state.

## Security boundary

At this checkpoint, `kryptasec scan`:

- does not send an HTTP request;
- does not crawl the target;
- does not execute shell commands;
- does not invoke an LLM;
- does not exploit vulnerabilities.

SQLite persistence stores scan metadata only. Active testing remains disabled in Phase 1.
