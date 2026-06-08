# Title

web_fetch Tool Specification

# High Level Description

This specification defines expected behavior for the web_fetch tool in agent/tool/web_fetch.
The tool fetches HTML from a URL and attempts to parse it through a language model.

# User Scenarios

1. As an orchestrating agent, I want to fetch webpage content so I can inspect source material.
2. As an orchestrating agent, I want parsed output from HTML so I can reason over textual content quickly.
3. As a maintainer, I want model and network failure handling to be explicit and testable.

# Functional Requirements

## FR-FETCH-001 URL Retrieval

- Apply must fetch response body from Request.Link.
- On success, raw body must be returned in Response.Html.
- If HTTP request fails, Apply must return an error.

## FR-FETCH-002 Model Parse Attempt

- Apply must submit fetched HTML to configured language model.
- If model returns valid agent message, Response.Parsed must equal returned message text.

## FR-FETCH-003 Parse Failure Handling

- If model call fails, Response.Parsed must contain failure reason text.
- If model result is empty, Response.Parsed must state empty result failure.
- If first model message type is not agent, Response.Parsed must state invalid response failure.

## FR-FETCH-004 Metadata And Schema

- Name must be web_fetch.
- RequestSchema and ResultSchema must be non-nil.
- DescribeRequest must include link.
- DescribeResult must return parsed text.

## FR-FETCH-005 Auto Policy

- Auto must return false.

# Non-Functional Requirements

1. Reliability: Expected network or model failures must return clear outcomes.
2. Testability: HTTP and model interactions must be mockable for local unit tests.
3. Determinism: Metadata methods must be stable for stable inputs.

# Definition of Done

1. FR-FETCH-001 to FR-FETCH-005 are covered by automated tests.
2. make test passes.
3. tasks.md is fully checked.
4. ammendments.md includes an update entry.

# Testing Methodology

1. Unit tests in agent/tool/web_fetch.
2. Use test HTTP server or mocked client behavior for retrieval paths.
3. Use mocked language model for success, error, empty, and invalid response cases.
4. Validate metadata, schema, and auto policy behavior.
