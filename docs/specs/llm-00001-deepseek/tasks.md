# Tasks

## Preparation

- [ ] Confirm DeepSeek adapter requirement IDs and wording with maintainers.
- [ ] Identify existing tests in agent/llm related to DeepSeek adapter behavior.

## Implementation

- [ ] Add startup and lifecycle tests for createDeepSeek/Start/Stop (FR-DEEPSEEK-001).
- [ ] Add Chat input mapping tests for system, user, tool request, and tool result messages (FR-DEEPSEEK-002).
- [ ] Add tests verifying allowedTools filter integration with ToolContainer.DeepSeekParams (FR-DEEPSEEK-003).
- [ ] Add completion tests for finish_reason=stop mapping to AGENT output with reasoning content (FR-DEEPSEEK-004).
- [ ] Add completion tests for finish_reason=tool_calls mapping to TOOL_REQUEST output and DescribeRequest integration (FR-DEEPSEEK-005).
- [ ] Add error-path tests for CreateChatCompletion failures returning ERROR messages and non-nil errors (FR-DEEPSEEK-006).
- [ ] Add helper conversion tests for DeepSeekToolFromJsonSchema (FR-DEEPSEEK-007).

## Validation

- [x] Run make test and fix failures.
- [ ] Update amendment/changelog artifacts per repo process.
- [ ] Mark completed tasks and document residual risks.
