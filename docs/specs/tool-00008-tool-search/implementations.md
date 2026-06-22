# Implementations

## Choices Made

- **Sync tool**: `tool_search` is a synchronous tool since search is a fast, in-memory operation with no async work.
- **Auto-approved**: Read-only search requires no user confirmation.
- **Keyword matching**: Initial implementation uses case-insensitive substring matching with weighted field scoring (name > description > schema parameters). Full-text or semantic matching is deferred.
- **Container reference**: The tool holds a reference to the `ToolContainer` interface, not a snapshot, so it always searches the live registry.

## Libraries Used

- Standard library: `strings` for keyword matching, `sort` for result ordering.
- Project libraries: `agent` for `SyncTool`/`ToolContainer` interfaces, `jsonschema` for schema generation, `format` for serialization.

## Implementation Preferences

1. Relevance scoring uses integer weights: name match = 3 points per keyword, description match = 2 points, schema match = 1 point.
2. Empty query returns empty slice, not an error.
3. `DescribeResult` formats results with one line per tool: `- <name>: <description> (auto: <yes/no>) [score: <n>]`.

## Caveats

1. Search is purely keyword-based — it does not understand synonyms or semantic similarity.
2. Tool descriptions and schemas are trusted as-is; no normalization or stemming is applied.
3. The `ToolContainer` must expose an enumeration method (`List()` or equivalent) — this is a prerequisite change to the container interface.
