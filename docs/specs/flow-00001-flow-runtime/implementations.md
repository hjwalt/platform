# Choices Made

- Preserved the original requirement intent and normalized identifiers to FR-FLOW-\* for template consistency.
- Kept acceptance criteria embedded in each functional requirement to preserve direct test mapping.
- Explicitly carried migration intent forward: new behavior belongs in `flow` and runtime, not legacy `flows` expansion.

# Libraries Used

- No new libraries required for this specification refactor.
- Existing repository test, logging, and metrics libraries remain the expected implementation basis.

# Implementation Preferences

- Prefer table-driven tests for runtime lifecycle, retry, and handler determinism checks.
- Keep key-scoped stateful scenarios explicit with per-key fixtures and concurrent test cases.
- Validate observability using structured log field assertions and metric label/value checks.

# Caveats

- This change refactors specification artifacts only; runtime/flow code behavior is unchanged until implementation tasks are completed.
- Final observability assertions may require alignment on canonical correlation metadata keys across components.
