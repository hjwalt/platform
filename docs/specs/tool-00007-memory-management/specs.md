# Title

Memory Management Tool Specification

# High Level Description

This specification defines a single tool for agent memory management:

- `memory`

The tool behavior is selected by a required request parameter:

- `operation = "get"`
- `operation = "update"`
- `operation = "clear"`

The tool operates on a memory store configured at setup time.
The managed canonical key is prefix-derived:

- no prefix: `memory`
- prefix `X`: `X`

The specification also requires a generic instantiation mechanism using a configurable prefix so multiple memory domains can coexist with isolated tool names and storage keys.

# User Scenarios

1. As an agent runtime integrator, I want to configure one memory tool so agents have one discoverable API surface for memory operations.
2. As an agent, I want to fetch current memory content from the canonical prefix-derived key so I can reason over persisted context.
3. As an agent, I want to update memory content deterministically so knowledge can be appended or replaced.
4. As an agent, I want to clear memory quickly so stale context can be removed without deleting tool configuration.
5. As a platform developer, I want to instantiate prefixed memory tool variants so multiple memory scopes (for example session, repo, user) can be used in the same runtime.

## Functional Requirements

## FR-TOOL-00007-001 Single Tool And Naming

- The memory tool surface must be exactly one tool name per instance: `memory`.
- The tool must be exposed through the standard tool interface used by the agent runtime.
- The request schema must include a required `operation` field with enum values: `get`, `update`, `clear`.

## FR-TOOL-00007-002 Canonical Key Derivation

- Tool configuration must provide the backing memory store and optional prefix.
- The effective canonical key must be:
  - `memory` when prefix is empty.
  - `<prefix>` when prefix is provided.
- Callers must not be allowed to override the derived canonical key directly.

## FR-TOOL-00007-003 `operation=get` Behavior

- The tool must read and return the full content of the derived canonical key.
- If the canonical key does not exist, the result must be deterministic and documented (for example, empty content with `exists=false`).
- `operation=get` must not mutate store state.

## FR-TOOL-00007-004 `operation=update` Behavior

- The tool must write content to the derived canonical key.
- Update mode must be explicit and deterministic (for example `replace` or `append`) and represented in request schema.
- `content` is required when `operation=update`.
- After successful update, a subsequent `operation=get` must return the updated content.

## FR-TOOL-00007-005 `operation=clear` Behavior

- The tool must remove logical content for the derived canonical key.
- Clear behavior must be deterministic and documented (for example delete key and/or write empty value).
- After successful clear, `operation=get` must return empty content according to the chosen existence semantics.

## FR-TOOL-00007-006 Prefix-Based Generic Instantiation

- The implementation must provide a generic constructor that accepts at minimum:
  - tool prefix
  - backing store configuration
- For a prefix `X`, the exposed tool name must be `X_memory`.
- Prefix handling must allow multiple independent instances to be registered simultaneously without name collisions.
- Prefix normalization and validation rules must be documented (for example allowed characters and separator policy).

## FR-TOOL-00007-007 Validation And Error Surface

- Invalid configuration (for example nil store or invalid prefix) must return actionable errors during setup.
- Runtime store operation failures must return stable, human-readable errors.
- Errors must distinguish configuration failures from runtime operation failures.

## FR-TOOL-00007-008 Metadata And Discoverability

- The tool must provide stable metadata (`Name`, request description, result description, and schemas).
- Metadata must include enough information for runtime UIs to explain operation-based behavior and prefix-derived canonical key semantics.

# Non-Functional Requirements

## NFR-TOOL-00007-001 Safety
- The tool must only operate against configured store/key scope.

## NFR-TOOL-00007-002 Predictability
- Read, write, and clear semantics must be deterministic across repeated runs.

## NFR-TOOL-00007-003 Testability
- Store behavior and prefix instantiation must be fully unit-testable without external services.

## NFR-TOOL-00007-004 Isolation
- Multiple prefixed instances must not interfere with each other.

## NFR-TOOL-00007-005 Simplicity
- Implementation should rely on Go standard library for filesystem operations where possible.

# Definition of Done

1. A specification-compliant implementation exists for one `memory` tool with operation-dispatch behavior.
2. Canonical key derivation semantics are enforced.
3. Prefix-based instantiation supports at least two concurrent memory domains in tests.
4. FR-TOOL-00007-001 through FR-TOOL-00007-008 are covered by automated tests.
5. `make test` passes without regressions.

# Testing Methodology

1. Unit tests for each `operation` mode with existing and missing key conditions.
2. Tests for operation request validation (`operation` required, unknown values rejected, `content` required for update).
3. Tests for update modes (`append` and `replace`) and read-after-write correctness.
4. Tests for clear semantics and idempotency.
5. Tests for invalid configuration and invalid prefix handling.
6. Multi-instance tests registering different prefixes and ensuring isolation.
7. Concurrency-focused tests for repeated updates and reads where practical.
