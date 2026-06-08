# Title

linux_shell Tool Specification

# High Level Description

This specification defines expected behavior for the linux_shell tool in agent/tool/linux_shell.
The tool executes shell commands from a configured base directory and returns command output.

# User Scenarios

1. As an orchestrating agent, I want to execute a shell command so that I can run local development tasks.
2. As a maintainer, I want execution failures surfaced as errors so that bad commands do not appear successful.
3. As an operator, I want readable request and result descriptions so tool activity is auditable.

# Functional Requirements

## FR-LINUX-001 Command Execution

- Apply must split Request.Command into executable and args.
- Apply must execute command in configured BaseDir.
- On success, combined output must be returned in Response.Result.

## FR-LINUX-002 Failure Propagation

- If command binary is missing, Apply must return an error.
- If command exits with non-zero status, Apply must return an error.

## FR-LINUX-003 Metadata And Schema

- Name must be linux_shell.
- Description must remain deterministic.
- RequestSchema and ResultSchema must be non-nil.
- DescribeRequest must include command text.
- DescribeResult must return response result text.

## FR-LINUX-004 Auto Policy

- Auto must return false.

# Non-Functional Requirements

1. Reliability: Expected error paths must return errors without panic.
2. Determinism: Metadata methods must return stable values for stable inputs.
3. Testability: Behavior must be verifiable with unit tests that run locally.

# Definition of Done

1. FR-LINUX-001 to FR-LINUX-004 are covered by automated tests.
2. make test passes.
3. tasks.md is fully checked.
4. ammendments.md includes an update entry.

# Testing Methodology

1. Unit tests in agent/tool/linux_shell.
2. Test success path command execution.
3. Test missing command and non-zero exit failures.
4. Test metadata, schema, and auto policy behavior.
