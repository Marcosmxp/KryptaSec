# ADR-0002: Go Core and React/TypeScript Dashboard

- Status: Accepted
- Date: 2026-09-18

## Context

KryptaSec requires:

- a distributable CLI;
- cancellable concurrent orchestration;
- sandbox/process management;
- HTTP/WebSocket APIs;
- a responsive local dashboard;
- strong separation between security core and UI.

## Decision

Use:

- Go for the CLI, application core, orchestration, policy engine, sandbox control, persistence services and local API;
- React + TypeScript for the dashboard;
- SQLite for initial local persistence;
- Docker as the initial sandbox runtime;
- Playwright for browser automation;
- provider adapters behind a KryptaSec-owned LLM interface.

## Consequences

Advantages:

- simple native CLI distribution;
- strong concurrency/cancellation primitives;
- low runtime overhead;
- clear backend/frontend boundary;
- React/TypeScript ecosystem for complex visualization.

Costs:

- two primary language ecosystems;
- Playwright may require Node/runtime support in sandbox tooling;
- some security libraries may require adapters or subprocess integration.

## Alternatives considered

### Python core

Strong AI/security ecosystem, but not chosen for the independent core architecture.

### TypeScript-only

Would simplify language count but is less attractive for process/sandbox orchestration and single-binary CLI goals.

### Rust core

Strong isolation/performance characteristics but higher implementation complexity for the initial team size.
