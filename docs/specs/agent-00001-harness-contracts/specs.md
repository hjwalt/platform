# Title

Agent Harness Contracts Specification

# High Level Description

This specification defines required behavior for the agent harness prompt lifecycle, tool execution boundaries, skill parsing, model-call outcomes, and auditability.
It standardizes typed interfaces and failure handling so harness behavior is predictable, testable, and safe.

# User Scenarios

1. As an agent maintainer, I want tool calls to enforce typed contracts so invalid model-produced inputs are rejected safely.
2. As a runtime operator, I want tool failures isolated so one failing call does not corrupt subsequent harness execution.
3. As a developer, I want skill parsing to produce deterministic structured output and clear parser errors.
4. As an integrator, I want model calls to return normalized success, timeout, and provider failure outcomes.
5. As a reviewer, I want traceable execution logs so I can reconstruct interaction order for success and failure paths.

# Functional Requirements

## FR-AGENT-00001-001 Typed Tool Invocation

- Tool execution must enforce request and response typing contracts.
- Given valid tool input payload, harness execution must return typed output that is serializable.
- Given invalid input payload, harness validation must reject execution with structured error details.

## FR-AGENT-00001-002 Tool Failure Isolation

- A failing tool call must not corrupt overall harness state.
- Given one failing tool invocation in a multi-step sequence, subsequent steps must follow explicit policy instead of undefined behavior.
- Given a tool panic or transport error, the harness must contain the failure and convert it into an explicit failure result.

## FR-AGENT-00001-003 Skill Parsing Contract

- Skill metadata and content must parse consistently from expected formats.
- Given valid frontmatter and body sections, parser execution must return a structured skill object.
- Given malformed frontmatter delimiters, parser execution must return a descriptive parse error.

## FR-AGENT-00001-004 Model Call Robustness

- Model execution must expose clear outcomes for success, timeout, and provider failure.
- Given provider success response, model execution must return content and metadata in a normalized shape.
- Given provider timeout or transport failure, model execution must apply configured retry policy or return terminal error behavior explicitly.

## FR-AGENT-00001-005 Auditability

- Harness execution must capture enough information to reconstruct tool and model interaction order.
- Given a successful run, logs must include traceable entries for each major step.
- Given a failed run, logs must expose failure location and reason without requiring debug-only instrumentation.

# Non-Functional Requirements

## NFR-AGENT-00001-001 Safety
- Tool input derived from model output must be validated before runtime dispatch.

## NFR-AGENT-00001-002 Security
- Sensitive values must be omitted from logs or explicitly redacted.

## NFR-AGENT-00001-003 Reliability
- Panic and transport-level failures must be contained and represented as explicit failures.

## NFR-AGENT-00001-004 Observability
- Logs must allow execution-sequence reconstruction in both success and failure paths.

## NFR-AGENT-00001-005 Maintainability
- Requirement IDs and acceptance statements must remain stable and test-mapped.

# Definition of Done

1. FR-AGENT-00001-001 through FR-AGENT-00001-005 are covered by automated tests.
2. Tests cover both happy paths and failure paths defined in acceptance scenarios.
3. Unknown tool names are rejected before runtime dispatch.
4. Logging behavior is validated to ensure traceability and redaction requirements.
5. tasks.md is fully checked and ammendments.md records implementation updates.

# Testing Methodology

1. Unit tests in `agent/tool` validate request/response typing, schema validation, and invalid payload rejection.
2. Harness tests in `agent/harness` validate failure isolation and panic/transport containment behavior.
3. Skill parser tests in `agent/skill` validate successful frontmatter parsing and malformed delimiter failures.
4. LLM integration tests in `agent/llm` with mocked providers validate normalized success, timeout, and provider-failure behavior.
5. Logger-backed tests verify execution ordering, failure discoverability, and redaction handling.
