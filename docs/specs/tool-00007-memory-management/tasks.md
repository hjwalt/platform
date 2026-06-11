# Tasks

## Preparation

- [x] Confirm operation enum for unified tool request: `get`, `update`, `clear`.
- [x] Confirm request validation rules (`content` required for `update`, optional `mode` defaults).
- [x] Align prefix validation rules with unified tool naming constraints.

## Implementation

- [x] Converge memory API surface from three tools to one `memory` tool (FR-MEM-001).
- [x] Add operation-dispatch request schema using `operation` parameter (FR-MEM-001, FR-MEM-003..005).
- [x] Ensure canonical key derivation is prefix-based (`memory` or `<prefix>`) (FR-MEM-002).
- [x] Keep generic constructor for prefix-based instantiation and registration (`X_memory`) (FR-MEM-006).
- [x] Add configuration validation and explicit error mapping for setup/runtime failures (FR-MEM-007).
- [x] Add metadata and schema coverage for unified tool discoverability (FR-MEM-008).

## Validation

- [x] Add unit tests for request validation across all operations.
- [x] Add unit tests for get/update/clear behavior through one tool entrypoint.
- [x] Add tests for prefix collisions and multi-instance isolation.
- [x] Run `make test` and resolve regressions.
- [x] Update `ammendments.md` and mark completed tasks.
- [x] Keep store-backed wiring (`state.Store`) and validate canonical key semantics.
