# Title

DeepSeek LLM Adapter Specification

# High Level Description

This specification defines expected behavior for the DeepSeek language-model adapter in agent/llm/deepseek.go.
The adapter must implement the shared agent.LanguageModel lifecycle and chat contract while translating platform messages and tool-calls to/from the DeepSeek SDK.

# User Scenarios

1. As an agent runtime, I want to initialize a DeepSeek client from config so chat requests can run using configured model credentials.
2. As an orchestrator, I want platform messages converted to DeepSeek message structures so conversation context is preserved.
3. As a tool-aware workflow, I want model-emitted tool calls converted into platform TOOL_REQUEST messages so tools can be executed.
4. As an operator, I want failures surfaced as both a returned error and an ERROR message payload.

# Functional Requirements

## FR-LLM-00001-001 Model Construction And Startup

- createDeepSeek must map ModelConfig fields into the adapter instance.
- Start must create and store a DeepSeek client with configured secret.
- Stop must be safe to call and must not panic.

## FR-LLM-00001-002 Input Message Mapping

- Chat must map each platform message type to DeepSeek message payloads.
- SYSTEM must map via DeepSeekSystemMessage.
- USER must map via DeepSeekUserMessage.
- TOOL_REQUEST must map via DeepSeekToolRequestMessage.
- TOOL_RESULT must map via DeepSeekToolResultMessage.

## FR-LLM-00001-003 Allowed Tool Filtering

- Chat must pass allowedTools to ToolContainer.DeepSeekParams.
- The completion request must include only tool definitions resolved through that filter call.

## FR-LLM-00001-004 Stop Completion Handling

- For each choice with finish_reason=stop, Chat must emit an AGENT message.
- AGENT output must preserve context and include response content.
- AGENT output must include reasoning content from DeepSeek choice message.

## FR-LLM-00001-005 Tool Call Completion Handling

- For each choice with finish_reason=tool_calls, Chat must emit TOOL_REQUEST messages.
- Each TOOL_REQUEST must include tool id, name, and arguments from DeepSeek tool call payload.
- TOOL_REQUEST output must use ToolContainer.DescribeRequest when description succeeds.
- TOOL_REQUEST output must include reasoning content from DeepSeek choice message.

## FR-LLM-00001-006 Completion Error Handling

- If CreateChatCompletion fails, Chat must return a non-nil error.
- On failure, Chat must return at least one ERROR message including the original error string.
- The ERROR message must preserve the input context.

## FR-LLM-00001-007 Schema Conversion Helper

- DeepSeekToolFromJsonSchema must convert a jsonschema.Schema into DeepSeek function parameters.
- Returned DeepSeek tool definition must include type=function plus function name and description.

# Non-Functional Requirements

## NFR-LLM-00001-001 Reliability
- Chat must tolerate mixed histories and empty allowed tool lists without panic.

## NFR-LLM-00001-002 Determinism
- Conversion helper and message mappers must return stable payloads for equivalent inputs.

## NFR-LLM-00001-003 Maintainability
- Lifecycle, mapping, and conversion logic must remain independently testable.

# Definition of Done

1. FR-LLM-00001-001 through FR-LLM-00001-007 are covered by automated tests.
2. make test passes.
3. tasks.md is fully checked.
4. Any behavior changes are recorded in repository amendment/changelog process.

# Testing Methodology

1. Unit tests in agent/llm for lifecycle, mapping, completion interpretation, and helper conversion.
2. Use stubbed DeepSeek client behavior for success and failure paths.
3. Verify allowed tool filtering via ToolContainer.DeepSeekParams interaction assertions.
4. Validate reasoning content propagation in AGENT and TOOL_REQUEST outputs.
