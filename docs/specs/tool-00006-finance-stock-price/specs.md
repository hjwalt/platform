# Title

finance_stock_price Tool Specification

# High Level Description

This specification defines expected behavior for the finance_stock_price tool in agent/tool/finance_stock_price.
The tool retrieves stock price information for a ticker symbol using the configured provider client.

# User Scenarios

1. As an orchestrating agent, I want to request a stock symbol so I can obtain a current market price.
2. As a maintainer, I want provider failures reported explicitly so downstream flows can handle them safely.
3. As an operator, I want a readable result summary with symbol, price, and currency.

# Functional Requirements

## FR-STOCK-001 Quote Retrieval

- Apply must request provider data for the requested symbol.
- On successful response with values, Apply must return Symbol, Currency, and first value in Response.Value.

## FR-STOCK-002 HTTP Status Handling

- If provider returns a non-200 status, Apply must return an error.
- On non-200 status, Response must echo request Symbol.

## FR-STOCK-003 Empty Data Handling

- If result values are missing, Apply must return an error.
- On missing values, Response must echo request Symbol.

## FR-STOCK-004 Metadata And Schema

- Name must be finance_stock_price.
- RequestSchema and ResultSchema must be non-nil.
- DescribeRequest must include symbol.
- DescribeResult must include symbol and formatted value.

## FR-STOCK-005 Auto Policy

- Auto must return false.

# Non-Functional Requirements

1. Reliability: Provider failure paths must return explicit errors.
2. Testability: Provider behavior must be validated with mocks/fakes.
3. Determinism: Metadata methods must return stable values.

# Definition of Done

1. FR-STOCK-001 to FR-STOCK-005 are covered by automated tests.
2. make test passes.
3. tasks.md is fully checked.
4. ammendments.md includes an update entry.

# Testing Methodology

1. Unit tests in agent/tool/finance_stock_price.
2. Use mocked provider responses for success, non-200, and empty data.
3. Validate metadata, schema, and auto policy behavior.
