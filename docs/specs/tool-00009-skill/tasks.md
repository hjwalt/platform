# Tasks

## Preparation

- [x] Confirm skill requirement IDs and wording with maintainers.
- [x] Identify skill registry wiring point in `configuration/tool.go`.
- [x] Verify skill markdown files exist under `docs/agents/` for integration testing.

## Implementation

- [x] Add `Request` and `Response` structs with JSON and jsonschema tags (FR-SKILL-001).
- [x] Implement `Create` constructor accepting a skill registry (FR-SKILL-007).
- [x] Implement `Apply` with lookup, success response, and missing-skill error (FR-SKILL-002, FR-SKILL-003, FR-SKILL-004).
- [x] Implement `Name`, `Description`, `RequestFormat`, `ResultFormat`, `RequestSchema`, `ResultSchema` (FR-SKILL-006).
- [x] Implement `DescribeRequest` and `DescribeResult` (FR-SKILL-003, FR-SKILL-006).
- [x] Implement `Auto` returning true (FR-SKILL-005).
- [x] Implement `AddToContainer` registration function (FR-SKILL-008).
- [x] Register tool in `configuration/tool.go` wiring.
- [x] Add or update tests for successful skill lookup (FR-SKILL-002, FR-SKILL-003).
- [x] Add or update tests for missing skill handling (FR-SKILL-004).
- [x] Add metadata and schema tests (FR-SKILL-006).
- [x] Add auto policy test for Auto=true (FR-SKILL-005).
- [x] Add empty registry test (FR-SKILL-004, FR-SKILL-007).

## Validation

- [x] Run `make test` and fix failures.
- [x] Update `ammendments.md` with implemented changes and date.
- [x] Mark completed tasks and note residual risks.

## Residual Risks

- The self-improving-agent skill is loaded into the registry but not registered as a subagent (only researcher-agent is). This is intentional — the skill tool makes skills available for on-demand context injection without spawning separate agents.
- Skill registry is currently built from hardcoded directory paths. Future work could make skill directories configurable or auto-discovered.
