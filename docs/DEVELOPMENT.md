# Development Guide

## 1. Branch workflow

```text
development
    |
    +-- feature/<scope>
    +-- fix/<scope>
    +-- refactor/<scope>
    +-- docs/<scope>

feature/fix/docs
    -> Pull Request
    -> development

development
    -> release PR
    -> main
```

Do not implement normal work directly on `main`. Prefer isolated feature branches even when working alone.

## 2. Commit convention

Use Conventional Commits:

```text
feat: add scan coordinator
fix: enforce redirect scope
docs: update sandbox architecture
refactor: isolate llm provider adapters
test: add policy engine fixtures
build: add Go CI workflow
chore: update repository metadata
```

## 3. Definition of done

A change is not done until applicable items are satisfied:

- implementation is complete;
- unit/integration tests pass;
- lint/type/static checks pass;
- security-sensitive behavior has negative tests;
- architecture/API documentation is updated;
- an ADR exists for significant design choices;
- project status/roadmap is updated when a milestone changes;
- new dependencies have a documented license/provenance check;
- no credentials, target data or secrets are committed.

## 4. Documentation maintenance

Every PR must answer:

1. Does this change architecture?
2. Does this change security boundaries?
3. Does this change CLI/API/config/report schemas?
4. Does this change a roadmap milestone?
5. Does this add a third-party dependency?

Any "yes" requires the corresponding documentation update in the same PR.

## 5. Planned local toolchain

Core:

- Go current supported stable release chosen when bootstrapping the module;
- `go test`;
- `go vet`;
- a project-selected Go linter pinned in CI.

Dashboard:

- TypeScript strict mode;
- React;
- a single package manager selected and pinned at bootstrap;
- lint, typecheck and unit tests.

Repository:

- Docker for integration fixtures;
- GitHub Actions for CI;
- generated artifacts never committed unless explicitly documented.

Version numbers should be pinned in repository files rather than relying on a developer machine.

## 6. Testing strategy

### Unit tests

For:

- policy decisions;
- target normalization;
- scan state machine;
- provider parsing;
- finding transitions;
- report serialization.

### Integration tests

Use intentionally vulnerable local fixtures. Never depend on arbitrary public internet targets in CI.

Integration fixtures must:

- run locally;
- have deterministic credentials/data;
- be disposable;
- document which vulnerability they model.

### End-to-end tests

A future E2E fixture should exercise:

```text
target -> discovery -> candidate -> validation -> evidence -> patch -> retest
```

## 7. Secrets

Never commit:

- LLM API keys;
- GitHub tokens;
- real user credentials;
- production cookies;
- real pentest artifacts containing sensitive target data.

Configuration should resolve secrets from environment variables or OS-backed secret storage when implemented.

## 8. Dependency policy

Before adding a dependency:

- verify active maintenance;
- verify license compatibility;
- prefer narrow, replaceable libraries;
- avoid dependencies that duplicate small standard-library capabilities;
- record strategically important choices in an ADR.

## 9. Review expectations

Security-critical modules require particularly strict review:

- `internal/policy`
- `internal/runtime`
- `internal/tools` network/shell capabilities
- secret handling
- patch/write operations

A successful build is not sufficient evidence that these modules are safe.
