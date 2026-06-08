# Tasks

## Preparation

- [ ] Confirm FR-FLOW-001 through FR-FLOW-006 wording with flow maintainers.
- [ ] Inventory existing tests in `runtime`, `flow/stateless`, and `flow/stateful` against requirement coverage.
- [ ] Confirm expected retry policy defaults and non-retriable classification sources.
- [ ] Confirm required success/failure telemetry fields (component ID, correlation metadata, retry state).

## Implementation

- [ ] Add or update runtime lifecycle tests for start-once transition and wait blocking semantics (FR-FLOW-001).
- [ ] Add or update shutdown tests for cancellation fan-out, worker termination, and cleanup hook execution (FR-FLOW-002).
- [ ] Add or update stateless tests for deterministic repeated execution on equivalent input/config (FR-FLOW-003).
- [ ] Add or update invalid-input stateless tests to assert explicit errors and no partial side effects (FR-FLOW-003).
- [ ] Add or update stateful tests for repeated-delivery logical correctness under at-least-once assumptions (FR-FLOW-004).
- [ ] Add or update stateful concurrency tests to assert key-scoped isolation (FR-FLOW-004).
- [ ] Add or update runtime retry tests for delay policy and max retry boundary behavior (FR-FLOW-005).
- [ ] Add or update retry tests to assert immediate surfacing for non-retriable failures (FR-FLOW-005).
- [ ] Add or update logger/metric assertions for success path metadata requirements (FR-FLOW-006).
- [ ] Add or update logger/metric assertions for failure reason and retry-state requirements (FR-FLOW-006).

## Validation

- [x] Run `make test` and resolve failures tied to FR-FLOW-001 through FR-FLOW-006.
- [ ] Validate requirement-to-test mapping coverage and identify any untested acceptance paths.
- [ ] Update `ammendments.md` with implementation progress and decisions.
- [ ] Mark completed checklist items and document residual risks.
