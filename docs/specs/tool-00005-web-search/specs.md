# Title

web_search Tool Specification

# High Level Description

This specification defines expected behavior for the web_search tool in agent/tool/web_search.
The tool performs a web query and maps provider results into structured response records.

# User Scenarios

1. As an orchestrating agent, I want to search with a term so I can discover relevant source URLs.
2. As a maintainer, I want clear errors and empty fallback results when provider calls fail.
3. As an operator, I want readable formatted result descriptions for logs and debugging.

# Functional Requirements

## FR-TOOL-00005-001 Query Execution

- Apply must execute a web search using Request.Term.
- On provider success, Apply must return mapped results in Response.Results.

## FR-TOOL-00005-002 Error Handling

- If provider call fails, Apply must return an error.
- On provider failure, Response.Results must be an empty array.

## FR-TOOL-00005-003 Result Mapping

- Each mapped result must include title, url, description, language, content_type, and extra_snippets from provider data.

## FR-TOOL-00005-004 Result Description Formatting

- DescribeResult must format each result with readable sections.
- DescribeResult must include snippets as individual list-style lines.

## FR-TOOL-00005-005 Metadata And Schema

- Name must be web_search.
- RequestSchema and ResultSchema must be non-nil.
- DescribeRequest must include search term.

## FR-TOOL-00005-006 Auto Policy

- Auto must return true.

# Non-Functional Requirements

## NFR-TOOL-00005-001 Reliability
- Provider errors must surface with explicit error returns.

## NFR-TOOL-00005-002 Testability
- Provider behavior must be mockable for unit tests.

## NFR-TOOL-00005-003 Determinism
- Metadata methods must return stable values.

# Definition of Done

1. FR-TOOL-00005-001 to FR-TOOL-00005-006 are covered by automated tests.
2. make test passes.
3. tasks.md is fully checked.
4. ammendments.md includes an update entry.

# Testing Methodology

1. Unit tests in agent/tool/web_search.
2. Use mocked provider client for success and failure paths.
3. Validate mapping fidelity and formatted description content.
4. Validate metadata, schema, and auto policy behavior.
