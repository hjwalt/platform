# Title

finance_fx_price Tool Specification

# High Level Description

This specification defines expected behavior for the finance_fx_price tool in agent/tool/finance_fx_price.
The tool retrieves FX price information for a currency pair using the configured provider client.

# User Scenarios

1. As an orchestrating agent, I want to request a base and quote pair so I can obtain a conversion value.
2. As a maintainer, I want provider failures reported explicitly so downstream flows can handle them safely.
3. As an operator, I want a readable result summary with currency pair and value.

# Functional Requirements

## FR-TOOL-00002-001 Quote Retrieval

- Apply must request provider data using pair format C:<Base><Quote>.
- On successful response with values, Apply must return Base, Quote, and first value in Response.Value.

## FR-TOOL-00002-002 HTTP Status Handling

- If provider returns a non-200 status, Apply must return an error.
- On non-200 status, Response must echo request Base and Quote.

## FR-TOOL-00002-003 Empty Data Handling

- If result values are missing, Apply must return an error.
- On missing values, Response must echo request Base and Quote.

## FR-TOOL-00002-004 Metadata And Schema

- Name must be finance_fx_price.
- RequestSchema and ResultSchema must be non-nil.
- DescribeRequest must include base and quote pair.
- DescribeResult must include pair and formatted value.

## FR-TOOL-00002-005 Auto Policy

- Auto must return false.

# Non-Functional Requirements

## NFR-TOOL-00002-001 Reliability
- Provider failure paths must return explicit errors.

## NFR-TOOL-00002-002 Testability
- Provider behavior must be validated with mocks/fakes.

## NFR-TOOL-00002-003 Determinism
- Metadata methods must return stable values.

# Definition of Done

1. FR-TOOL-00002-001 to FR-TOOL-00002-005 are covered by automated tests.
2. make test passes.
3. tasks.md is fully checked.
4. ammendments.md includes an update entry.

# Testing Methodology

1. Unit tests in agent/tool/finance_fx_price.
2. Use mocked provider responses for success, non-200, and empty data.
3. Validate metadata, schema, and auto policy behavior.
