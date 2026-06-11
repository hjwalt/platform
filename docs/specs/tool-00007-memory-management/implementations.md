# Choices Made

- Converged the memory API surface to one sync tool: `memory`.
- Introduced a required request field `operation` with values `get`, `update`, and `clear`.
- Kept storage delegated to `state.Store` to decouple behavior from filesystem implementation.
- Kept canonical key derivation by prefix:
  - `memory` when prefix is empty.
  - `<prefix>` when prefix is provided.
- Kept explicit update modes for `operation=update`: `replace` (default) and `append`.
- Kept clear behavior deterministic via deletion or equivalent logical emptying on the canonical key.
- Converged prefixed naming to one tool per prefix: `prefix_memory`.
- Continued using `memory_tool.Configuration` with `Store state.Store` and `Prefix string`.
- Continued using `configuration.MemoryConfiguration` (`BaseDir`, `Prefix`) for file-backed wiring in the configuration layer.

# Libraries Used

- No third-party libraries were added.
- Implementation uses Go standard library packages: `regexp`, `sync`, and `time`.
- Storage is via `github.com/hjwalt/platform/state` interface.
- In-memory store for tests: `github.com/hjwalt/platform/state/memory`.

# Implementation Preferences

- Keep operation dispatch centralized in `agent/tool/memory` to minimize discoverability surface for agents.
- Canonical key is `MemoryKey` (`"memory"`) for empty prefix and `<prefix>` for prefixed instances.
- Prefix is validated with `^[a-zA-Z0-9]+(?:_[a-zA-Z0-9]+)*$`.
- Shared mutex is used inside the tool instance to serialize get/update/clear store operations.
- Results include `Key` (the store key), operation echo, and operation metadata (`exists`, `mode`, `bytes`, `cleared`).

# Caveats

- Request validation must reject unknown `operation` values and reject `update` when `content` is missing.
- Backward compatibility may require temporary aliases for `memory_get`, `memory_update`, and `memory_clear` if external callers still rely on the old names.
- Prefix naming policy must remain compatible with existing tool-registration constraints.
- `operation=get` should clearly document `exists=false` behavior for missing keys.
