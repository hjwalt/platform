# Tasks

## Preparation

- [x] Confirm final update semantics for `memory_update` (replace-only vs replace+append mode).
- [x] Confirm clear semantics for `memory_clear` (truncate vs delete+recreate).
- [x] Align prefix validation rules with existing tool naming constraints.

## Implementation

- [x] Add implementation for `memory_get` reading `<root_path>/memory.md` (FR-MEM-003).
- [x] Add implementation for `memory_update` writing `<root_path>/memory.md` with deterministic mode handling (FR-MEM-004).
- [x] Add implementation for `memory_clear` applying deterministic clear behavior (FR-MEM-005).
- [x] Add generic constructor for prefix-based tool instantiation and registration (FR-MEM-006).
- [x] Update canonical file naming so prefixed instances use `<prefix>.md` and non-prefixed instances use `memory.md`.
- [x] Add configuration validation and explicit error mapping for setup/runtime failures (FR-MEM-007).
- [x] Add metadata and schema coverage for all tools (FR-MEM-008).

## Validation

- [x] Add unit tests for missing file, existing file, update, and clear flows.
- [x] Add tests for prefix collisions and multi-instance isolation.
- [x] Run `make test` and resolve regressions.
- [x] Update `ammendments.md` and mark completed tasks.
- [x] Migrate storage from filesystem to `state.Store` interface.
- [x] Update `memory_get`, `memory_update`, `memory_clear` sub-tools to accept `state.Store` and key.
- [x] Update `memory_tool.Configuration` to accept `state.Store` instead of `BaseDir`.
- [x] Add `MemoryConfiguration` to `configuration` package for file-backed wiring.
- [x] Update `configuration/tool.go` to create file store from `MemoryConfiguration.BaseDir`.
- [x] Update `main.go` to use `configuration.MemoryConfiguration`.
- [x] Update tests to use in-memory store instead of `t.TempDir()`.
