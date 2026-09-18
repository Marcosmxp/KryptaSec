# Contributing to KryptaSec

## Branch policy

The repository uses the following workflow:

- `main` contains stable, release-ready code.
- `development` is the integration branch for active development.
- New work should use `feature/<name>`, `fix/<name>`, or `hotfix/<name>`.
- Feature and fix branches must target `development`.
- `development` reaches `main` only through a release pull request.

## External contributions

External contributors should fork the repository, create a branch in their fork, and open a pull request against `development`.

Do not target `main` directly for normal development.
