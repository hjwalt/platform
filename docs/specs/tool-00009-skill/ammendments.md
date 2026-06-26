# Ammendments

## 2026-06-25 — Initial specification

- Created specification for the `skill` tool (tool-00009).
- Defined 8 functional requirements covering tool interface, skill lookup, response content, error handling, auto policy, metadata, registry access, and registration.

## 2026-06-25 — Implementation

- Implemented `skill` SyncTool at `agent/tool/skill/tool.go`.
- Tool accepts a skill registry (`map[string]agent_skill.Skill`) at construction time.
- Case-insensitive exact-match lookup with O(1) performance.
- Missing/empty skill names return `Found: false` response (no Go error, no panic).
- `Auto()` returns `true` (read-only, no side effects).
- `DescribeResult` includes full markdown body for context injection into the LLM conversation.
- Wired into `configuration/tool.go` via `buildSkillRegistry()` helper that scans `./docs/agents/` directories.
- Currently registers two skills: `researcher-agent` and `self-improving-agent`.
- 14 unit tests covering all FRs (FR-SKILL-001 through FR-SKILL-008).
- Full test suite passes with no regressions.
