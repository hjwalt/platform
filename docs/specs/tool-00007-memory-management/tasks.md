# Tasks

## Preparation

- [ ] Confirm final update semantics for `memory_update` (replace-only vs replace+append mode).
- [ ] Confirm clear semantics for `memory_clear` (truncate vs delete+recreate).
- [ ] Align prefix validation rules with existing tool naming constraints.

## Implementation

- [ ] Add implementation for `memory_get` reading `<root_path>/memory.md` (FR-MEM-003).
- [ ] Add implementation for `memory_update` writing `<root_path>/memory.md` with deterministic mode handling (FR-MEM-004).
- [ ] Add implementation for `memory_clear` applying deterministic clear behavior (FR-MEM-005).
- [ ] Add generic constructor for prefix-based tool instantiation and registration (FR-MEM-006).
- [ ] Add configuration validation and explicit error mapping for setup/runtime failures (FR-MEM-007).
- [ ] Add metadata and schema coverage for all tools (FR-MEM-008).

## Validation

- [ ] Add unit tests for missing file, existing file, update, and clear flows.
- [ ] Add tests for prefix collisions and multi-instance isolation.
- [ ] Run `make test` and resolve regressions.
- [ ] Update `ammendments.md` and mark completed tasks.
