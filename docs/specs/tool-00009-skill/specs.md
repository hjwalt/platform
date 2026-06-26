# Title

Skill Tool Specification

# High Level Description

This specification defines the `skill` tool — a synchronous context-injection tool that loads an on-demand skill (markdown playbook) into the current agent's conversation context. The tool accepts a `name` string identifying the skill and returns the skill's body, description, and allowed tools as a tool result. The returned body becomes part of the conversation, giving the model the guidance it needs for subsequent steps.

The `skill` tool implements **progressive disclosure**: skills are only loaded when the model determines it needs their guidance, rather than being pre-loaded into the system prompt. This keeps the base context lean while making specialized playbooks available on demand.

Unlike the `subagent` tool (which spawns a new asynchronous agent execution), the `skill` tool is synchronous and returns content directly to the calling agent — the skill's guidance runs in the same agent context, not a child context.

# User Scenarios

1. As an agent, I want to load a skill by name (e.g., `"code-review"`) when I'm about to review code, so I get the review playbook's guidance injected into my context right when I need it.
2. As an agent, I want the skill's body, description, and allowed-tool hints returned so I can follow the playbook's instructions for subsequent tool calls.
3. As an agent, I want a clear error when a skill name is not found, so I can ask the user for clarification or try an alternative.
4. As a platform developer, I want skills to be registered in a name-indexed registry so lookup is deterministic and fast.
5. As an operator, I want skill loading to be read-only and auto-approved so it doesn't add permission friction.

# Functional Requirements

## FR-SKILL-001 Tool Naming And Interface

- The tool must be named `skill`.
- The tool must be exposed through the standard `SyncTool` interface used by the agent runtime.
- The request schema must include a required `name` field of type `string`.
- The result schema must include `name` (string), `description` (string), `body` (string), and `allowed_tools` (array of strings).

## FR-SKILL-002 Skill Lookup

- The tool must look up skills by exact name match in a skill registry provided at construction time.
- Name matching must be case-insensitive after trimming whitespace.
- The skill registry must be a map-like structure keyed by skill name.
- Lookup must be O(1) — no linear scan or fuzzy matching.

## FR-SKILL-003 Response Content

- On a successful lookup, the result must include the skill's `name`, `description`, `body` (the full markdown playbook), and `allowed_tools`.
- `body` must contain the complete markdown content from the skill's source file (the portion after the YAML frontmatter).
- `allowed_tools` must reflect the skill's declared allowed tools from its frontmatter.
- `DescribeResult` must format the response to clearly separate the skill description from the body content.

## FR-SKILL-004 Missing Skill Handling

- When a skill name is not found in the registry, the tool must return an error result (not panic).
- The error must include the requested name so the agent can report it to the user.
- An empty or whitespace-only `name` must be treated as not-found.

## FR-SKILL-005 Auto Policy

- `Auto()` must return `true` — skill loading is read-only with no side effects.

## FR-SKILL-006 Metadata And Schemas

- `Name()` must return `"skill"`.
- `RequestSchema()` and `ResultSchema()` must be non-nil and valid JSON Schema.
- `RequestFormat()` and `ResultFormat()` must use the standard format pipeline (`format.Json[T]()`).
- `DescribeRequest` must include the skill name being requested.

## FR-SKILL-007 Skill Registry Access

- The tool must receive a skill registry (map of skill name to `agent_skill.Skill`) at construction time.
- The tool must not mutate the registry or any registered skills.
- The tool must handle an empty registry gracefully (return a not-found error for any lookup).

## FR-SKILL-008 Registration

- An `AddToContainer` function must register the tool into a `ToolContainer` as a sync tool.
- A `Create` constructor must accept a skill registry and return a `SyncTool[Request, Response]`.

# Non-Functional Requirements

1. Safety: The tool must be read-only — it must not modify skill registrations, configuration, or state.
2. Performance: Lookup must complete in O(1) time with no external service calls or file I/O (skills are pre-loaded at startup).
3. Determinism: Identical lookups against an unchanged registry must return identical results.
4. Testability: The tool must accept a skill registry as a constructor parameter, enabling fully mocked unit tests.
5. Memory: The tool must not retain or cache request data beyond the scope of a single `Apply` call.

# Definition of Done

1. A specification-compliant implementation exists at `agent/tool/skill/`.
2. The tool is registered via `AddToContainer()` in `configuration/tool.go`.
3. FR-SKILL-001 through FR-SKILL-008 are covered by automated tests.
4. `make test` passes without regressions.
5. `tasks.md` is fully checked.
6. `ammendments.md` includes an initial entry.

# Testing Methodology

1. Unit tests for successful skill lookup returning correct name, description, body, and allowed_tools.
2. Unit tests for missing skill (name not in registry) returning a descriptive error.
3. Unit tests for empty string and whitespace-only name inputs.
4. Unit tests for case-insensitive name matching.
5. Unit tests for metadata (Name, RequestSchema, ResultSchema, Auto, DescribeRequest, DescribeResult).
6. Unit tests for empty registry handling.
7. Integration test verifying skill lookup with at least 2 registered skills.
