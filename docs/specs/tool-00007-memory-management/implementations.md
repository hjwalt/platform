# Choices Made

- Implemented all memory operations in `agent/tool/memory` as three independent sync tools backed by shared configuration.
- Storage is now delegated to a `state.Store` interface, decoupling tools from the filesystem.
- The canonical key in the store is derived from prefix:
  - `memory` when prefix is empty.
  - `<prefix>` when prefix is provided.
- Added explicit update modes for `memory_update`: `replace` (default) and `append`.
- Implemented clear behavior as `store.Delete` on the canonical key.
- Added generic prefix-aware constructor so names are `prefix_memory_get`, `prefix_memory_update`, and `prefix_memory_clear`.
- `memory_tool.Configuration` now accepts `Store state.Store` and `Prefix string` (removed `BaseDir`).
- Added `configuration.MemoryConfiguration` (`BaseDir`, `Prefix`) for file-backed wiring in the configuration layer.
- `configuration/tool.go` creates a `file_store` from `MemoryConfiguration.BaseDir` and passes it to `memory_tool.Configuration`.

# Libraries Used

- No third-party libraries were added.
- Implementation uses Go standard library packages: `regexp`, `sync`, and `time`.
- Storage is via `github.com/hjwalt/platform/state` interface.
- In-memory store for tests: `github.com/hjwalt/platform/state/memory`.

# Implementation Preferences

- Keep each memory operation implementation in its own tool folder:
  - `agent/tool/memory_get`
  - `agent/tool/memory_update`
  - `agent/tool/memory_clear`
    while using `agent/tool/memory` as the composition and registration wrapper.
- Canonical key is `MemoryKey` (`"memory"`) for empty prefix and `<prefix>` for prefixed instances.
- Prefix is validated with `^[a-zA-Z0-9]+(?:_[a-zA-Z0-9]+)*$`.
- Shared mutex is used inside the tool set instance to serialize get/update/clear store operations.
- Results include `Key` (the store key) and operation metadata (`exists`, `mode`, `bytes`, `cleared`).

# Caveats

- After `memory_clear` (which calls `store.Delete`), `memory_get` returns `exists=false` with empty content.
  This is a slight behavior change from the filesystem implementation where `exists=true` with empty content was returned after clear.
- Prefix naming policy must remain compatible with existing tool-registration constraints.
- `memory_get` returns `exists=false` with empty content when the store returns zero-length bytes for the key.
- File-backed store (`state/file`) requires the directory to exist; the trailing `/` in `BaseDir + "/"` is required for correct path construction.
