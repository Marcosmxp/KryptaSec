# KryptaSec

Open-source autonomous security testing for modern applications.

## Branching model

- `main`: stable/release branch. No day-to-day development.
- `development`: integration branch for active development.
- `feature/*`: feature work, merged into `development` through pull requests.
- `fix/*`: bug fixes, merged into `development` through pull requests.
- `hotfix/*`: urgent production fixes, reviewed before reaching `main`.

## Contribution flow

1. Fork the repository or create an authorized feature branch.
2. Branch from `development`.
3. Open a pull request targeting `development`.
4. After validation, releases are promoted from `development` to `main`.

Direct development on `main` is intentionally avoided.
