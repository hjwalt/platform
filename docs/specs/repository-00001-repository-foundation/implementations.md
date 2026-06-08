# Choices Made

- Migrated REQ-FOUNDATION-_ requirements into template-aligned FR-FOUNDATION-_ functional requirements.
- Preserved original acceptance intent while normalizing structure to SDD folder conventions.
- Kept manual runtime baseline validation explicit due to local service dependency requirements.

# Libraries Used

- No new libraries required for this specification refactor.
- Existing repository tooling via `make` commands remains the validation baseline.

# Implementation Preferences

- Prefer command-driven validation in CI and local workflows for reproducibility.
- Keep requirement IDs visible in tests, review notes, and pull request descriptions.
- Favor deterministic checks that avoid environment-specific side effects outside declared build artifacts.

# Caveats

- This migration changes specification artifacts only; implementation and tests must be completed via the task checklist.
- Integration and runtime baseline validation depend on local service availability and environment setup.
