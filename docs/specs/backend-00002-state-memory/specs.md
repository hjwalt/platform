# Title

In-Memory State Store Specification

# High Level Description

This specification defines a new in-memory implementation of `state.Store` located in `state/memory` with package name `memory_store`.
The store keeps state in process memory and is safe for simultaneous access by multiple goroutines performing reads, writes, and key enumeration.

# User Scenarios

1. As a flow developer, I want a `state.Store` implementation that keeps data in memory for fast local state operations.
2. As a runtime integrator, I want concurrent goroutines to safely update and read the same store without races or panics.
3. As a caller, I want missing-key reads to return an empty `state.State` value rather than a not-found failure.
4. As a test author, I want deterministic behavior for `Read`, `Write`, and `Keys` so unit tests are stable and predictable.

## Functional Requirements

## FR-STATE-MEM-001 Store Placement And Packaging

- The implementation must live under the folder `state/memory`.
- The package declaration must be `package memory_store`.
- The exported store type must satisfy the `state.Store` interface from `state/types.go`.

## FR-STATE-MEM-002 In-Memory Persistence Semantics

- The store must hold values in memory only and must not persist state to disk or external systems.
- Given `Write(ctx, state.State{Id: X, Value: V})`, when the call succeeds, then subsequent `Read(ctx, X)` returns `Value` equal to `V`.
- Given multiple writes for the same `Id`, when `Read` is called, then the latest successfully written value for that `Id` is returned.
- Given a write with an empty `Id`, when `Write` is called, then the operation returns an explicit validation error.

## FR-STATE-MEM-003 Concurrent Safety

- The store must be safe for simultaneous `Read`, `Write`, and `Keys` calls from multiple goroutines.
- Given concurrent writes to distinct keys, when operations complete, then all successfully written keys are present and readable.
- Given concurrent writes to the same key, when operations complete, then no race, panic, or partial byte corruption occurs.
- Given concurrent `Keys` and `Write` operations, when `Keys` returns, then it returns a valid snapshot and does not panic.

## FR-STATE-MEM-004 Read And Keys Behavior

- Given an existing key, when `Read` is called, then the returned state has matching `Id`, stored `Value`, and a populated `Timestamp`.
- Given a missing key, when `Read` is called, then it returns `state.State{Id: requestedId, Value: []byte{}}` and no not-found error.
- When `Keys` is called, it must return all known key identifiers in deterministic order.
- Returned key slices and value bytes must not expose mutable internal storage that can corrupt store state when modified by callers.

## FR-STATE-MEM-005 Lifecycle Behavior

- `Start()` must initialize internal structures needed for operation and be safe to call once during runtime startup.
- `Stop()` must complete without panic and release in-memory references so the store can be garbage-collected.
- Lifecycle behavior must be documented for caller expectations (for example, whether operations before `Start` are allowed).

## FR-STATE-MEM-006 Error Surface

- The store must return stable, actionable errors for invalid inputs and internal misuse.
- Validation failures (for example empty key) must be distinguishable from operational failures.
- Concurrent access must not produce data races under `go test -race`.

# Non-Functional Requirements

1. Concurrency safety: No data races or unsafe memory access under concurrent test workloads.
2. Performance: Common operations (`Read`, `Write`) should remain amortized constant time for key lookup and update.
3. Predictability: `Keys` output ordering must be deterministic for repeatable tests.
4. Isolation: No filesystem, network, or external dependency is required.
5. Testability: Behavior must be fully verifiable with unit tests, including race-enabled test runs.

# Definition of Done

1. A store implementation exists in `state/memory` with package name `memory_store`.
2. The implementation satisfies `state.Store` and compiles with the repository.
3. Functional requirements FR-STATE-MEM-001 through FR-STATE-MEM-006 are covered by automated tests.
4. Concurrency tests pass under race detection for simultaneous reads/writes/keys operations.
5. Missing-key read, deterministic key ordering, and copy-on-read/copy-on-keys behavior are validated.
6. `make test` passes without regressions.

# Testing Methodology

1. Contract tests: verify `Read`, `Write`, `Keys`, `Start`, and `Stop` match expected `state.Store` semantics.
2. Missing-key tests: validate no-error empty-state behavior for unknown keys.
3. Concurrency tests: run high-parallel read/write/keys scenarios and assert correctness.
4. Race tests: execute package tests with race detection to ensure no data races.
5. Determinism tests: assert stable ordering from `Keys` across repeated calls.
6. Defensive-copy tests: mutate returned slices and verify store internals remain unchanged.
