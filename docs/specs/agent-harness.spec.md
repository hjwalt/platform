# Agent Harness Specification

Status: Proposed
Owner: Agent Maintainers
Last Updated: 2026-06-08

## Scope

Behavior of agent prompt handling, tool execution, skill parsing, and model interactions.

## Goals

- Ensure predictable tool-call lifecycle and typed boundaries.
- Define failure handling for model and tool integration points.
- Preserve safety and auditability in agent execution.

## Non-Goals

- Provider-specific model quality benchmarks.
- UI-level chat presentation behavior.

## Requirements

### REQ-AGENT-001: Typed Tool Invocation

Tool execution must enforce request and response typing contracts.

Acceptance scenarios:

1. Given valid tool input payload, when the harness executes a tool, then typed output is returned and serializable.
2. Given invalid input payload, when the harness validates request, then execution is rejected with structured error details.

Test mapping:

- agent tool package unit tests.

### REQ-AGENT-002: Tool Failure Isolation

A failing tool call must not corrupt overall harness state.

Acceptance scenarios:

1. Given one failing tool invocation in a multi-step sequence, when execution continues, then subsequent steps are governed by explicit policy rather than undefined behavior.
2. Given tool panic or transport error, when harness handles failure, then panic is contained and converted into explicit failure result.

Test mapping:

- agent harness failure-path tests.

### REQ-AGENT-003: Skill Parsing Contract

Skill metadata and content must parse consistently from expected formats.

Acceptance scenarios:

1. Given valid frontmatter and body sections, when parser runs, then structured skill object is produced.
2. Given malformed frontmatter delimiters, when parser runs, then descriptive parse error is returned.

Test mapping:

- agent skill parser tests.

### REQ-AGENT-004: Model Call Robustness

Model execution must expose clear outcomes for success, timeout, and provider failure.

Acceptance scenarios:

1. Given provider success response, when model call completes, then content and metadata are returned in normalized shape.
2. Given provider timeout or transport failure, when model call completes, then retry or terminal error behavior follows configured policy.

Test mapping:

- agent llm integration tests with mocked providers.

### REQ-AGENT-005: Auditability

Harness execution must capture enough information to reconstruct tool and model interaction order.

Acceptance scenarios:

1. Given a successful run, when execution logs are inspected, then each major step has a traceable entry.
2. Given a failed run, when logs are inspected, then failure location and reason are discoverable without debug-only instrumentation.

Test mapping:

- logger-backed harness tests.

## Security And Safety Constraints

- Tool input from model output must be validated before execution.
- Sensitive values must not be emitted in logs unless explicitly redacted.
- Unknown tool names must be rejected before runtime dispatch.
