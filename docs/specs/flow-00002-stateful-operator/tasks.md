# Tasks

## Preparation

- [ ] Review FR-STATE-001 through FR-STATE-009 wording with flow maintainers.
- [ ] Inventory the existing `flow/stateful` package (no test files present) to confirm zero coverage baseline.
- [ ] Confirm expected behavior for the output+error precedence rule (output over error) with maintainers.
- [ ] Verify that `flow.Store`, `flow.Producer`, `flow.Message`, `flow.State`, and `flow.ExtractMetadata` types are stable and won't change under this spec cycle.
- [ ] Confirm `either.Either` and `optional.Optional` APIs are stable.

## Implementation

- [ ] Write fake `flow.Store[ST]` implementation for testing (in-memory map with call recording).
- [ ] Write fake `flow.Producer[V]` implementation for testing (slice capture with call recording).
- [ ] Write constructor test: `NewOperator` returns `flow.Handler[IV]` with all fields populated (FR-STATE-001).
- [ ] Write key extraction error test: `StateKey` failure returns error, no store read or further processing (FR-STATE-003).
- [ ] Write state read error test: `Store.Read` failure returns error, no state update or handler (FR-STATE-005).
- [ ] Write state update right-path (error short-circuit) test: `StateUpdate` returns `either.Right`, error message produced, no state written, no handler called (FR-STATE-004).
- [ ] Write state write error test: `Store.Write` failure returns error, no handler execution (FR-STATE-005).
- [ ] Write full success path test: key → read → update (Left) → write → handler → output message produced with correct metadata (FR-STATE-002, FR-STATE-006, FR-STATE-007).
- [ ] Write handler error path test: handler returns present error, error message produced with correct metadata (FR-STATE-006, FR-STATE-007).
- [ ] Write filter/sink path test: handler returns absent output and absent error, no messages produced, nil returned (FR-STATE-006).
- [ ] Write output-over-error precedence test: handler returns both present output and present error, only output message produced (FR-STATE-006).
- [ ] Write context propagation test: operator name is set in log context (FR-STATE-008).
- [ ] Write metadata propagation test: input metadata flows through extractor to output/error messages (FR-STATE-007).
- [ ] Write type parameter smoke test: instantiate operator with distinct IV, OV, ST, ERR types to validate generic constraints (FR-STATE-001, FR-STATE-009).
- [ ] Write state-update-key-isolation test: distinct keys produce independent state read/write calls.

## Validation

- [ ] Run `make test` and resolve failures in `flow/stateful/`.
- [ ] Validate requirement-to-test mapping: every FR-STATE-001 through FR-STATE-009 must have at least one test case.
- [ ] Verify test coverage for all branches in `Handle`: key error, read error, update right, write error, handler (output, error, both, neither), producer errors.
- [ ] Update `ammendments.md` with implementation progress and any decisions made.
- [ ] Mark completed checklist items and document residual risks.
