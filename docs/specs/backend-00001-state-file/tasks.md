# Tasks

## Preparation

- [ ] Confirm the intended `state/file` contract for missing reads, path readiness, and lifecycle expectations.
- [ ] Inventory current `state/file` test coverage and identify gaps against FR-STATE-FILE-001 through FR-STATE-FILE-005.
- [ ] Decide whether `Start` should validate or create the configured directory before runtime usage.

## Implementation

- [ ] Add or update write/read persistence coverage for `.dat`-backed state values (FR-STATE-FILE-001, FR-STATE-FILE-002).
- [ ] Add or update missing-key read coverage to preserve empty-state semantics without not-found errors (FR-STATE-FILE-002).
- [ ] Add or update `Keys` coverage for suffix trimming, directory exclusion, and directory-read failures (FR-STATE-FILE-003).
- [ ] Implement or document lifecycle behavior for configured-path readiness and shutdown semantics (FR-STATE-FILE-004).
- [ ] Add or update tests for sentinel-wrapped read and write failures (FR-STATE-FILE-005).

## Validation

- [ ] Run focused `state/file` tests and confirm FR-STATE-FILE-001 through FR-STATE-FILE-005 coverage.
- [ ] Run `make test` and confirm repository-wide regression safety.
- [ ] Update ammendments.md with implementation outcomes and mark completed tasks.
