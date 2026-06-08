# Choices Made

- Used the `backend` specification type because the memory state store is an internal backend component.
- Scoped requirements directly to the existing `state.Store` interface to keep compatibility explicit.
- Included explicit concurrent-safety requirements because the feature must support simultaneous goroutine updates.
- Implemented synchronization with `sync.RWMutex` and an internal map for constant-time in-memory lookups and writes.
- Implemented `Start` as idempotent initialization and `Stop` as state teardown (map release).

# Libraries Used

- No new libraries are required by this specification artifact.
- Recommended implementation can use the Go standard library (`sync`, `sort`, `time`) only.

# Implementation Preferences

- Prefer `sync.RWMutex` or equivalent synchronization to protect shared in-memory maps.
- Prefer pairing `Lock`/`RLock` with `defer Unlock`/`defer RUnlock` to ensure unlock always happens on all return paths.
- Return defensive copies for byte slices and key slices to prevent caller mutation of internal state.
- Keep key ordering deterministic (for example, sorted ascending) to simplify tests and reduce flaky behavior.

# Caveats

- State is process-local only and is intentionally lost across process restarts.
- `Stop` clears all in-memory state as part of lifecycle teardown.
