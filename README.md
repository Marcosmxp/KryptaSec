# KryptaSec

KryptaSec is a free, source-available autonomous application-security platform focused on software built with modern frameworks and AI-assisted development.

> Status: Phase 1 core CLI is in development. No production security scanner has been released yet.

## Project principles

- Independent implementation: KryptaSec is not a source-code port of Strix or any other product.
- Security by authorization: active testing must remain inside an explicitly approved scope.
- Evidence before findings: vulnerabilities should be reported as confirmed only when supported by reproducible evidence.
- Local-first: the initial product runs on the user's machine and uses the user's selected LLM provider.
- Provider-neutral: cloud and local LLMs are accessed through a dedicated abstraction.
- Documentation is part of the product: architecture and decisions must be updated with code changes.

## Planned stack

- Core/CLI/orchestration: Go
- Dashboard: React + TypeScript
- Local API: Go HTTP/WebSocket API
- Storage: SQLite
- Sandbox: Docker first, Podman-compatible later
- Browser automation: Playwright
- Reports: JSON, SARIF, Markdown and HTML
- AI providers: OpenAI-compatible APIs, Anthropic, Gemini, OpenRouter and local models through adapters

## Current development CLI

The Phase 1 branch currently provides:

```bash
kryptasec --version
kryptasec doctor
kryptasec scan ./local-project
kryptasec scan --scope-host staging.example.com https://staging.example.com
kryptasec scan status <scan-id>
```

Scan jobs are persisted locally in SQLite under the KryptaSec data directory. Every state transition uses an expected-state transactional update before the next state is accepted.

The current `scan` command still performs no active testing or remote network request.

See [docs/CLI.md](docs/CLI.md) for the exact current contract.

## Branch model

- `main`: stable releases only
- `development`: integration branch
- `feature/*`: new functionality
- `fix/*`: bug fixes
- `hotfix/*`: urgent release fixes

Normal work starts from `development` and returns through pull requests.

## Documentation

Start at [docs/README.md](docs/README.md).

Key documents:

- [Architecture](docs/ARCHITECTURE.md)
- [Roadmap](docs/ROADMAP.md)
- [Development guide](docs/DEVELOPMENT.md)
- [CLI](docs/CLI.md)
- [Dependencies](docs/DEPENDENCIES.md)
- [Security model](docs/SECURITY_MODEL.md)
- [Legal and provenance policy](docs/LEGAL_AND_PROVENANCE.md)
- [Project status](docs/PROJECT_STATUS.md)
- [Architecture decisions](docs/adr/)

## Legal note

Functional similarity to another security product does not make a rewrite legally independent by itself. KryptaSec follows a clean-room policy: no copying of source code, prompts, UI assets, documentation, tests, internal names or other protected expression from reference products. Third-party dependencies and borrowed code, when intentionally used, must be tracked with their licenses and attribution requirements.

The final distribution license has not yet been selected and must be decided before the first public release.
