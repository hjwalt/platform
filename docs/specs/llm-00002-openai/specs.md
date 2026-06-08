# Title

OpenAI LLM Adapter Specification

# High Level Description

This specification defines expected behavior for the OpenAI language-model adapter in agent/llm/openai.go.
The adapter must implement the shared agent.LanguageModel lifecycle and chat contract while translating platform messages and tool-calls to/from the OpenAI SDK.

# User Scenarios

1. As an agent runtime, I want to initialize the OpenAI client from config so chat requests can be executed against a configured endpoint.
2. As an orchestrator, I want platform messages converted to OpenAI message formats so existing conversation history is preserved.
3. As a tool-aware workflow, I want model tool-calls converted back into platform tool requests so tools can be executed.
4. As an operator, I want transport or API failures surfaced as both an error return and an ERROR message payload.

# Functional Requirements

## FR-OPENAI-001 Model Construction And Startup

- createOpenAi must map ModelConfig fields into the adapter instance.
- Start must create an OpenAI client using endpoint and secret from config.
- Stop must be safe to call and must not panic.

## FR-OPENAI-002 Input Message Mapping

- Chat must map each platform message type to OpenAI message params.
- SYSTEM must map via OpenAiSystemMessage.
- USER must map via OpenAiUserMessage.
- TOOL_REQUEST must map via OpenAiToolRequestMessage.
- TOOL_RESULT must map via OpenAiToolResultMessage.

## FR-OPENAI-003 Allowed Tool Filtering

- Chat must pass allowedTools to ToolContainer.OpenAiParamsFiltered.
- The completion request must include only filtered tool definitions.

## FR-OPENAI-004 Stop Completion Handling

- For each choice with finish_reason=stop, Chat must emit an AGENT message.
- Emitted AGENT messages must preserve input context and include assistant text content.

## FR-OPENAI-005 Tool Call Completion Handling

- For each choice with finish_reason=tool_calls, Chat must emit TOOL_REQUEST messages.
- Each TOOL_REQUEST must include tool id, name, and arguments from OpenAI tool call payload.
- TOOL_REQUEST message body must be generated through ToolContainer.DescribeRequest when description succeeds.

## FR-OPENAI-006 Completion Error Handling

- If the OpenAI completion request fails, Chat must return a non-nil error.
- On failure, Chat must return at least one ERROR message including the original error string.
- The ERROR message must preserve the input context.

## FR-OPENAI-007 Schema Conversion Helpers

- OpenAiToolSchema must derive JSON schema from generic request type M and forward it through OpenAiFromJsonSchema.
- OpenAiFromJsonSchema must convert schema payloads into OpenAI function parameter format.
- Returned OpenAI tool definitions must include function name and description.

# Non-Functional Requirements

1. Reliability: Chat must tolerate mixed message histories and empty tool lists without panics.
2. Determinism: Mapping helpers must produce stable output for equivalent inputs.
3. Maintainability: Message and schema conversion functions must remain independently unit-testable.

# Definition of Done

1. FR-OPENAI-001 through FR-OPENAI-007 are covered by automated tests.
2. make test passes.
3. tasks.md is fully checked.
4. Any behavior changes are recorded in repository amendment/changelog process.

# Testing Methodology

1. Unit tests in agent/llm for startup, mapping, and completion output behavior.
2. Use mocked or stubbed OpenAI completion transport/client behavior for success and failure paths.
3. Verify tool filtering by asserting calls through ToolContainer.OpenAiParamsFiltered.
4. Validate schema helper outputs for representative request structs.
