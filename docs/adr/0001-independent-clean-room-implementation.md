# ADR-0001: Independent Clean-Room Implementation

- Status: Accepted
- Date: 2026-09-18

## Context

KryptaSec is intended to provide autonomous application-security capabilities in a market that includes Strix and other security platforms.

A simple rewrite of another project's source into a different language would not provide the desired implementation independence and could preserve protected expression.

## Decision

KryptaSec will use a clean-room-oriented engineering policy.

Product capabilities may be informed by public standards, security research, public product behavior and independently written requirements. Implementation must be written from KryptaSec specifications rather than translated/copied source.

When third-party open-source code is deliberately reused, it must be treated as an explicit licensed dependency or attributed source, not disguised as independent implementation.

## Consequences

Positive:

- clearer provenance;
- independent architecture;
- easier license auditing;
- reduced coupling to another project's internals.

Cost:

- functionality must be designed and implemented independently;
- some capabilities will take longer than a direct fork;
- provenance/documentation discipline is mandatory.

## Alternatives considered

### Fork Strix

Not selected because the project goal is independent architecture and ownership of the implementation.

### Translate/rewrite Strix into Go

Rejected. Language translation does not make copied implementation independent.
