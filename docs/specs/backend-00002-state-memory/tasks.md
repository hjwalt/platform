# Tasks

## Preparation

- [x] Review `state.Store` contract and align `state/memory` expectations for read, write, keys, start, and stop behavior.
- [x] Confirm package naming and path constraints: folder `state/memory`, package `memory_store`.
- [x] Define synchronization strategy for safe multi-goroutine access.

## Implementation

- [x] Create `state/memory` package and implement a store that satisfies `state.Store` (FR-STATE-MEM-001).
- [x] Implement write/read semantics for latest-value behavior and missing-key handling (FR-STATE-MEM-002, FR-STATE-MEM-004).
- [x] Add synchronization to make `Read`, `Write`, and `Keys` safe under concurrent goroutine access (FR-STATE-MEM-003).
- [x] Implement deterministic `Keys` ordering and defensive copy behavior for returned keys/values (FR-STATE-MEM-004).
- [x] Implement lifecycle behavior for `Start` and `Stop` with clear caller expectations (FR-STATE-MEM-005).
- [x] Add stable validation and operational error handling (FR-STATE-MEM-006).

## Validation

- [x] Add/extend unit tests for contract behavior and missing-key semantics.
- [x] Add/extend concurrency and race-focused tests for simultaneous updates and reads.
- [x] Run focused package tests for `state/memory`.
- [x] Run `make test` and ensure no regressions.
- [x] Update ammendments.md and mark completed tasks after implementation.
