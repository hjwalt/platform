# Title

subagent Tool Specification

# High Level Description

This specification defines expected behavior for the subagent async tool in agent/tool/subagent.
The tool delegates work by producing a start message for a sub-agent with constrained context.

# User Scenarios

1. As an orchestrating agent, I want to delegate a prompt to a subagent so a sub-agent can process it asynchronously.
2. As a maintainer, I want delegated contexts to preserve parent linkage and allowed tool constraints.
3. As an operator, I want stable metadata for predictable registration and tracing.

# Functional Requirements

## FR-SUBAGENT-001 Async Delegation

- Send must produce exactly one agent start message per call.
- Produced message must include toolCall.Id.

## FR-SUBAGENT-002 Delegated Context Propagation

- Delegated context must carry ParentContext from parent agent context.
- Delegated context must set SystemMessage from configured subagent prompt.
- Delegated context must include configured AllowedTools.

## FR-SUBAGENT-003 Constructor Behavior

- Create must configure Name, Description, Subagent body, and AllowedTools from Configuration.
- FromSubagent must configure Name, Description, Subagent body, and AllowedTools from subagent definition.

## FR-SUBAGENT-004 Metadata And Schema

- RequestSchema must be non-nil.
- Name and Description must be deterministic for configured values.
- DescribeRequest must include tool name and prompt content.

## FR-SUBAGENT-005 Auto Policy

- Auto must return true.

# Non-Functional Requirements

1. Reliability: Send must not panic on valid producer invocation paths.
2. Testability: Producer interactions must be verifiable with fakes or mocks.
3. Determinism: Metadata output must remain stable for stable input.

# Definition of Done

1. FR-SUBAGENT-001 to FR-SUBAGENT-005 are covered by automated tests.
2. make test passes.
3. tasks.md is fully checked.
4. ammendments.md includes an update entry.

# Testing Methodology

1. Unit tests in agent/tool/subagent.
2. Use fake producer to assert message count and message payload.
3. Validate context propagation, metadata, schema, and auto policy.
