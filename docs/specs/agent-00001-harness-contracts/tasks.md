# Tasks

## Preparation

- [ ] Confirm FR-AGENT-001 through FR-AGENT-005 wording with agent maintainers.
- [ ] Inventory existing tests in `agent/tool`, `agent/harness`, `agent/skill`, and `agent/llm`.
- [ ] Define explicit failure policy for multi-step tool sequences.

## Implementation

- [ ] Add or update typed request/response contract tests for tool invocation (FR-AGENT-001).
- [ ] Add invalid payload validation and structured error contract tests (FR-AGENT-001).
- [ ] Add harness failure-isolation tests for single-step failure in multi-step runs (FR-AGENT-002).
- [ ] Add panic and transport error containment tests with explicit failure result mapping (FR-AGENT-002).
- [ ] Add parser tests for valid frontmatter/body parsing to structured skill objects (FR-AGENT-003).
- [ ] Add parser tests for malformed frontmatter delimiters and descriptive parse errors (FR-AGENT-003).
- [ ] Add model-call tests for normalized success shape including metadata (FR-AGENT-004).
- [ ] Add timeout and provider-failure tests for retry/terminal behavior policy (FR-AGENT-004).
- [ ] Add logger-backed tests for traceable major-step ordering and failure discoverability (FR-AGENT-005).
- [ ] Add redaction-focused tests to ensure sensitive values are not emitted in logs (FR-AGENT-005).
- [ ] Add unknown-tool rejection tests before dispatch.

## Validation

- [x] Run `make test` and resolve test failures.
- [ ] Validate requirement-to-test mapping coverage for FR-AGENT-001 through FR-AGENT-005.
- [ ] Update ammendments.md with completed implementation details.
- [ ] Mark completed checklist items and document residual risks.
