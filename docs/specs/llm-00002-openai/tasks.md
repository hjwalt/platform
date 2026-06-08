# Tasks

## Preparation

- [ ] Confirm OpenAI adapter requirement IDs and wording with maintainers.
- [ ] Identify existing tests in agent/llm related to OpenAI adapter behavior.

## Implementation

- [ ] Add startup and lifecycle tests for createOpenAi/Start/Stop (FR-OPENAI-001).
- [ ] Add Chat input message mapping tests for system, user, tool request, and tool result (FR-OPENAI-002).
- [ ] Add tests verifying allowedTools filter integration with ToolContainer.OpenAiParamsFiltered (FR-OPENAI-003).
- [ ] Add completion tests for finish_reason=stop mapping to AGENT output (FR-OPENAI-004).
- [ ] Add completion tests for finish_reason=tool_calls mapping to TOOL_REQUEST output and DescribeRequest integration (FR-OPENAI-005).
- [ ] Add error-path tests for completion failures returning ERROR messages and non-nil errors (FR-OPENAI-006).
- [ ] Add helper conversion tests for OpenAiToolSchema and OpenAiFromJsonSchema (FR-OPENAI-007).

## Validation

- [x] Run make test and fix failures.
- [ ] Update amendment/changelog artifacts per repo process.
- [ ] Mark completed tasks and document residual risks.
