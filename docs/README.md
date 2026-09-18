# KryptaSec Documentation

This directory is the source of truth for the technical and operational design of KryptaSec.

## Required documents

| Document | Purpose |
| --- | --- |
| [ARCHITECTURE.md](ARCHITECTURE.md) | System boundaries, components, data flow and repository layout |
| [ROADMAP.md](ROADMAP.md) | Milestones from foundation to v1.0 |
| [DEVELOPMENT.md](DEVELOPMENT.md) | Engineering workflow, quality gates and branch conventions |
| [SECURITY_MODEL.md](SECURITY_MODEL.md) | Authorization, sandboxing, scope and safety constraints |
| [LEGAL_AND_PROVENANCE.md](LEGAL_AND_PROVENANCE.md) | Clean-room and dependency provenance policy |
| [PROJECT_STATUS.md](PROJECT_STATUS.md) | Current phase, completed work and next tasks |
| [adr/](adr/) | Architecture Decision Records |

## Documentation rule

A pull request must update documentation when it changes any of the following:

- architecture or component boundaries;
- public CLI/API behavior;
- security or authorization rules;
- persistence schema;
- supported providers or integrations;
- build/deployment workflow;
- release process;
- project roadmap or implementation status.

If a change introduces a meaningful architectural tradeoff, create an ADR instead of only editing prose.
