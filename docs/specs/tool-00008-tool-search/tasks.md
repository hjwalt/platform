# Tasks

## Preparation

- [ ] Add `List() []ToolInfo` or equivalent enumeration method to `ToolContainer` interface
- [ ] Define `ToolInfo` struct with name, description, auto, and schema fields
- [ ] Implement container enumeration in `agent/util/container/container.go`
- [ ] Define `tool_search` request and response types

## Implementation

- [ ] Create `agent/tool/tool_search/` package
- [ ] Implement `tool_search` SyncTool with keyword matching and relevance scoring
- [ ] Implement `AddToContainer()` constructor accepting `ToolContainer` reference
- [ ] Register `tool_search` in `configuration/tool.go` via `RegisterTools()`
- [ ] Add `ToolSearch` configuration entry to `ToolConfiguration` in `configuration/types.go`

## Validation

- [ ] Write unit tests for query normalization and edge cases
- [ ] Write unit tests for keyword matching and relevance scoring
- [ ] Write unit tests for result ordering and empty-container behavior
- [ ] Write unit tests for metadata (Name, Schema, Auto, DescribeRequest, DescribeResult)
- [ ] Run `make test` and confirm no regressions
- [ ] Update spec index in `/docs/memory/spec-driven-development.md`
