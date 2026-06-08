# Choices Made

- Used `tool` specification type because the feature introduces runtime-invokable agent tools.
- Captured all three operations in one spec to keep shared path and prefix behavior consistent.
- Fixed canonical filename to `memory.md` as a hard requirement to simplify interoperability.
- Required generic prefix instantiation so multiple memory scopes can coexist in one runtime.

# Libraries Used

- No new libraries are required by this specification artifact.
- Recommended implementation can use Go standard library packages such as `os`, `path/filepath`, and `strings`.

# Implementation Preferences

- Prefer resolving canonical path with safe path join logic and explicit normalization.
- Prefer deterministic update semantics exposed via request schema rather than implicit behavior.
- Prefer returning structured result fields such as content length, existed flag, and resolved path when useful for diagnostics.
- Prefer small helper abstractions for file IO so unit tests can simulate IO failure paths.

# Caveats

- Prefix naming policy must remain compatible with existing tool-registration constraints.
- If future requirements demand multiple files, this spec must be amended because it currently enforces a single canonical filename.
