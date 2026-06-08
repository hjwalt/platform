# Title

Memory Management Tools Specification

# High Level Description

This specification defines three tools for agent memory management as markdown files:

- `memory_get`
- `memory_update`
- `memory_clear`

Each tool operates on a memory root path provided through tool configuration.
The managed memory file name is prefix-derived under the resolved root path:

- no prefix: `memory.md`
- prefix `X`: `X.md`

The specification also requires a generic instantiation mechanism using a configurable prefix so multiple memory domains can coexist with isolated tool names and storage roots.

# User Scenarios

1. As an agent runtime integrator, I want to configure a root memory path so memory tools consistently read and write one canonical file.
2. As an agent, I want to fetch current memory content from the canonical prefix-derived memory file so I can reason over persisted context.
3. As an agent, I want to update memory content deterministically so knowledge can be appended or replaced.
4. As an agent, I want to clear memory quickly so stale context can be removed without deleting tool configuration.
5. As a platform developer, I want to instantiate prefixed tool variants so multiple memory scopes (for example session, repo, user) can be used in the same runtime.

## Functional Requirements

## FR-MEM-001 Tool Set And Naming

- The tool family must include exactly these operations: `memory_get`, `memory_update`, and `memory_clear`.
- The tools must be exposed through the standard tool interface used by the agent runtime.
- Each tool must have request and result schemas suitable for structured invocation.

## FR-MEM-002 Root Path And Canonical File

- Tool configuration must require a root file path used as the memory base path.
- The effective memory document path must be:
  - `<root_path>/memory.md` when prefix is empty.
  - `<root_path>/<prefix>.md` when prefix is provided.
- Callers must not be allowed to override the derived memory file name directly.
- If the root path does not exist, behavior must be explicit and deterministic (either create during update/clear or return a clear error, as defined by implementation).

## FR-MEM-003 memory_get Behavior

- `memory_get` must read and return the full content of the derived canonical file.
- If the derived canonical file does not exist, the result must be deterministic and documented (for example, empty string content with exists=false).
- `memory_get` must not mutate filesystem state.

## FR-MEM-004 memory_update Behavior

- `memory_update` must write content to the derived canonical file.
- Update mode must be explicit and deterministic (for example, replace or append), and represented in request schema.
- After successful update, a subsequent `memory_get` must return the updated content.
- The operation must be atomic enough to avoid partial writes from process interruption where practical.

## FR-MEM-005 memory_clear Behavior

- `memory_clear` must remove logical content of the derived canonical file.
- Clear behavior must be deterministic and documented (truncate to empty file and/or recreate empty canonical file).
- After successful clear, `memory_get` must return empty content.

## FR-MEM-006 Prefix-Based Generic Instantiation

- The implementation must provide a generic constructor that accepts at minimum:
  - tool prefix
  - root path
- For a prefix `X`, exposed tool names must be:
  - `X_memory_get`
  - `X_memory_update`
  - `X_memory_clear`
- Prefix handling must allow multiple independent instances to be registered simultaneously without name collisions.
- Prefix normalization and validation rules must be documented (for example allowed characters and separator policy).

## FR-MEM-007 Validation And Error Surface

- Invalid configuration (for example empty root path or invalid prefix) must return actionable errors during setup.
- Runtime file operation failures must return stable, human-readable errors.
- Errors must distinguish configuration failures from runtime IO failures.

## FR-MEM-008 Metadata And Discoverability

- Each tool must provide stable metadata (`Name`, request description, result description, and schemas).
- Metadata must include enough information for runtime UIs to explain root path semantics and prefix-derived canonical filename behavior.

# Non-Functional Requirements

1. Safety: Tools must only operate within configured root path and canonical `memory.md` location.
2. Predictability: Read, write, and clear semantics must be deterministic across repeated runs.
3. Testability: File behavior and prefix instantiation must be fully unit-testable without external services.
4. Isolation: Multiple prefixed instances must not interfere with each other.
5. Simplicity: Implementation should rely on Go standard library for filesystem operations where possible.

# Definition of Done

1. A specification-compliant implementation exists for `memory_get`, `memory_update`, and `memory_clear`.
2. Root path configuration and fixed `memory.md` semantics are enforced.
3. Prefix-based instantiation supports at least two concurrent memory domains in tests.
4. FR-MEM-001 through FR-MEM-008 are covered by automated tests.
5. `make test` passes without regressions.

# Testing Methodology

1. Unit tests for each tool operation with existing and missing `memory.md` conditions.
2. File-system tests using temporary directories to validate canonical path behavior.
3. Tests for update modes (append/replace if supported) and read-after-write correctness.
4. Tests for clear semantics and idempotency.
5. Tests for invalid configuration and invalid prefix handling.
6. Multi-instance tests registering different prefixes against different roots.
7. Concurrency-focused tests for repeated updates and reads where practical.
