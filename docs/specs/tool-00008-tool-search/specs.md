# Title

Tool Search Tool Specification

# High Level Description

This specification defines the `tool_search` tool — a discovery tool that lets agents search for registered tools by describing what they need in natural language. The tool accepts a `query` string and returns matching tool references with names, descriptions, and capabilities. It enables agents to find tools when they aren't sure of the exact tool name, and supports broad queries to discover related tools in a single call.

The tool complements the `availableDeferredTools` mechanism by providing a programmatic search surface over the same tool registry.

# User Scenarios

1. As an agent, I want to search for tools by describing a capability (e.g., "I need to run shell commands") so I can discover the `linux_shell` tool without knowing its exact name.
2. As an agent, I want to use broad queries (e.g., "finance") to find all finance-related tools in one call rather than making multiple narrow searches.
3. As an agent, I want to discover tool parameter requirements through search results so I can form valid tool calls.
4. As a platform developer, I want the search to return stable, structured results so agent reasoning is deterministic.
5. As an operator, I want search to be read-only and auto-approved so it doesn't add permission friction.

# Functional Requirements

## FR-TOOL-00008-001 Tool Naming And Interface

- The tool must be named `tool_search`.
- The tool must be exposed through the standard `SyncTool` interface used by the agent runtime.
- The request schema must include a required `query` field of type `string`.
- The tool must implement `Auto() bool` returning `true` (read-only, no side effects).

## FR-TOOL-00008-002 Query Input

- `query` must accept natural language descriptions of desired tool capabilities.
- Empty or whitespace-only queries must return an empty result set (not an error).
- Query strings must be trimmed of leading/trailing whitespace before processing.

## FR-TOOL-00008-003 Search Scope

- The tool must search across all tools registered in the `ToolContainer` at invocation time.
- Searchable fields for each tool must include: `Name`, `Description`, and parameter names/types from `RequestSchema`.
- The search must cover both sync and async tools.

## FR-TOOL-00008-004 Matching Semantics

- Matching must be case-insensitive keyword/substring matching against searchable fields.
- Matches against `Name` must score higher than matches against `Description`.
- Matches against `Description` must score higher than matches against parameter schemas.
- Tools with no keyword matches must be excluded from results.

## FR-TOOL-00008-005 Result Structure

- Each result must include at minimum: `name` (string), `description` (string), and `auto` (bool).
- Each result must include a `relevance` score indicating match strength (higher = better match).
- Results must be returned in descending order of relevance.

## FR-TOOL-00008-006 Result Formatting

- `DescribeResult` must format results as a readable list with name, description, and auto-policy for each match.
- `DescribeRequest` must include the search query text.

## FR-TOOL-00008-007 Container Access

- The tool must receive a reference to the `ToolContainer` at construction time.
- The tool must not mutate the container or any registered tools.
- The tool must handle an empty container gracefully (return empty results).

## FR-TOOL-00008-008 Metadata And Schemas

- `Name()` must return `"tool_search"`.
- `RequestSchema()` and `ResultSchema()` must be non-nil and valid JSON Schema.
- `RequestFormat()` and `ResultFormat()` must use the standard format pipeline.
- Metadata must include enough information for runtime UIs to explain query semantics and result structure.

# Non-Functional Requirements

## NFR-TOOL-00008-001 Safety
- The tool must be read-only — it must not modify tool registrations, configuration, or state.

## NFR-TOOL-00008-002 Performance
- Search must complete in O(n) time relative to registered tool count, with no external service calls.

## NFR-TOOL-00008-003 Determinism
- Identical queries against an unchanged tool set must return identical results.

## NFR-TOOL-00008-004 Testability
- The tool must accept a `ToolContainer` interface, enabling fully mocked unit tests.

## NFR-TOOL-00008-005 Extensibility
- The search interface must accommodate future enhancements (semantic search, category filtering) without breaking the request/response schema.

# Definition of Done

1. A specification-compliant implementation exists at `agent/tool/tool_search/`.
2. The tool is registered via `AddToContainer()` in `configuration/tool.go`.
3. FR-TOOL-00008-001 through FR-TOOL-00008-008 are covered by automated tests.
4. `make test` passes without regressions.
5. `tasks.md` is fully checked.
6. `ammendments.md` includes an initial entry.

# Testing Methodology

1. Unit tests for query normalization (empty, whitespace-only, trimming).
2. Unit tests for matching semantics with single and multi-keyword queries.
3. Unit tests for result ordering by relevance score.
4. Unit tests for empty container and no-match scenarios.
5. Unit tests for case-insensitive matching.
6. Unit tests for metadata (Name, RequestSchema, ResultSchema, Auto, DescribeRequest, DescribeResult).
7. Integration test verifying tool discovery with at least 3 registered tools of mixed types (sync + async).
