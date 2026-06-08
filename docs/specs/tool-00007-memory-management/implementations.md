# Choices Made

- Implemented all memory operations in `agent/tool/memory` as three independent sync tools backed by shared configuration.
- Fixed canonical filename derivation to prefix-aware behavior and resolved operations through:
  - `<root_path>/memory.md` when prefix is empty.
  - `<root_path>/<prefix>.md` when prefix is provided.
- Added explicit update modes for `memory_update`: `replace` (default) and `append`.
- Implemented clear behavior as deterministic truncate-to-empty by atomically writing an empty file.
- Added generic prefix-aware constructor so names are `prefix_memory_get`, `prefix_memory_update`, and `prefix_memory_clear`.
- Wired optional registration from configuration using `ToolConfiguration.Memory []memory_tool.Configuration`.

# Libraries Used

- No third-party libraries were added.
- Implementation uses Go standard library packages: `os`, `path/filepath`, `regexp`, `strings`, and `sync`.

# Implementation Preferences

- Keep each memory operation implementation in its own tool folder:
  - `agent/tool/memory_get`
  - `agent/tool/memory_update`
  - `agent/tool/memory_clear`
    while using `agent/tool/memory` as the composition and registration wrapper.
- Canonical path is built with `filepath.Join(baseDir, derivedFileName)`.
- `derivedFileName` is `memory.md` for empty prefix and `<prefix>.md` for prefixed instances.
- Prefix is validated with `^[a-zA-Z0-9]+(?:_[a-zA-Z0-9]+)*$`.
- File writes use atomic temp-file and rename strategy to reduce partial-write risk.
- Shared mutex is used inside the tool set instance to serialize get/update/clear file operations.
- Results include resolved path and operation metadata (`exists`, `mode`, `bytes`, `cleared`).

# Caveats

- Prefix naming policy must remain compatible with existing tool-registration constraints.
- `memory_get` returns `exists=false` with empty content when the derived canonical file is missing.
- If future requirements demand multiple files, this spec must be amended because it currently enforces a single canonical filename.
