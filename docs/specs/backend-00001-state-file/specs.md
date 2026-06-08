# Title

File State Store Specification

# High Level Description

This specification defines the backend contract for `state/file`, a filesystem-backed implementation of `state.Store`.
It covers how state values are persisted, retrieved, enumerated, and initialized when the store is used by platform components that need lightweight local durability without an external database.

# User Scenarios

1. As a developer, I want to persist state values to the local filesystem so components can recover data between process runs.
2. As a caller, I want reading a missing key to return an empty state object instead of a hard failure so first-write flows remain simple.
3. As an operator, I want store initialization failures to be explicit so misconfigured filesystem paths are diagnosable.
4. As an integrator, I want key enumeration to return all stored state identifiers so higher-level flows can inspect or replay persisted state.

## Functional Requirements

## FR-STATE-FILE-001 File-Backed Persistence

- The store must persist each `state.State` value as a file under the configured path.
- Given a state object with an `Id` and `Value`, when `Write` is called, then the store writes the value bytes to a file named `<Id>.dat` under the configured path.
- Given a successful write, when the same `Id` is later read, then the returned state contains the previously written bytes.

## FR-STATE-FILE-002 Read Semantics

- The store must support deterministic read behavior for existing and missing keys.
- Given an existing persisted file, when `Read` is called for its identifier, then the store returns a `state.State` with the requested `Id`, the file contents in `Value`, and a populated `Timestamp`.
- Given no persisted file exists for an identifier, when `Read` is called, then the store returns a `state.State` with the requested `Id`, an empty byte slice, and no not-found error.
- Given the filesystem returns a non-not-found read error, when `Read` is called, then the error is wrapped with the store's read failure sentinel.

## FR-STATE-FILE-003 Key Enumeration

- The store must enumerate persisted keys from the configured directory.
- Given the configured path contains files ending in `.dat`, when `Keys` is called, then the returned identifiers are the file names with the `.dat` suffix removed.
- Given directory entries include subdirectories, when `Keys` is called, then subdirectories are excluded from the returned identifiers.
- Given directory enumeration fails, when `Keys` is called, then the error is returned to the caller.

## FR-STATE-FILE-004 Store Lifecycle And Configuration

- The store must expose lifecycle behavior compatible with the `state.Store` interface.
- Given a valid filesystem path configuration, when `Start` is called, then the store completes successfully and the path is ready for subsequent read, write, and key operations.
- Given the store is no longer needed, when `Stop` is called, then shutdown completes without panicking and without corrupting persisted files.
- Given the configured path is invalid or inaccessible, when lifecycle or I/O operations depend on that path, then failures are explicit and actionable.

## FR-STATE-FILE-005 Error Surface

- The store must provide stable error semantics for callers.
- Given a write operation fails, when `Write` returns, then the error is wrapped with the store's write failure sentinel.
- Given a read operation fails for reasons other than file absence, when `Read` returns, then the error is wrapped with the store's read failure sentinel.
- Error messages and sentinels must remain specific enough for tests and callers to distinguish read and write failures.

# Non-Functional Requirements

1. Simplicity: The implementation should remain lightweight and avoid external service dependencies.
2. Durability: Successfully written values must survive process restarts as long as the backing filesystem remains intact.
3. Diagnosability: Filesystem and configuration failures must be visible through explicit returned errors.
4. Predictability: File naming and key enumeration rules must be stable across supported environments.
5. Testability: Core behavior must be verifiable with isolated temporary-directory tests.

# Definition of Done

1. FR-STATE-FILE-001 through FR-STATE-FILE-005 are covered by automated tests or explicit manual validation steps.
2. The store writes and reads `.dat` files using the configured directory with no unexpected file naming deviations.
3. Missing-key reads return an empty `state.State` without a not-found error.
4. Key enumeration returns persisted identifiers and excludes directories.
5. Read and write failures expose stable sentinel-wrapped errors.
6. Lifecycle behavior is documented and validated for configured-path readiness.

# Testing Methodology

1. Persistence validation: write a state value, read it back, and verify byte-for-byte equality (FR-STATE-FILE-001, FR-STATE-FILE-002).
2. Missing-key validation: read an identifier with no backing file and verify empty value semantics without error (FR-STATE-FILE-002).
3. Enumeration validation: create `.dat` files and subdirectories in a temporary path, run `Keys`, and verify only file-backed identifiers are returned (FR-STATE-FILE-003).
4. Error validation: use inaccessible or invalid paths to verify read and write sentinel wrapping behavior (FR-STATE-FILE-004, FR-STATE-FILE-005).
5. Lifecycle validation: exercise `Start` and `Stop` around temporary-directory setup and confirm no panic or silent misconfiguration (FR-STATE-FILE-004).
