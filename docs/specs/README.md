# Repository Specs For Spec-Driven Development

This directory defines executable, test-first specifications for the platform repository.

## How To Use

1. Pick a target spec file.
2. Select one requirement ID not yet implemented.
3. Write or update failing tests first.
4. Implement only enough code to make tests pass.
5. Run make test and update the requirement status.

## Spec Status Legend

- Proposed: drafted but not started.
- In Progress: tests and implementation are actively being developed.
- Implemented: acceptance tests pass.
- Verified: implemented and explicitly covered in regression tests.

## Spec Index

- repository-foundation.spec.md: Repository-wide engineering constraints and delivery gates.
- flow-runtime.spec.md: Dataflow and runtime lifecycle behavior.
- agent-harness.spec.md: Agent tool-call, model execution, and safety behavior.
- tool-00001-linux-shell/: linux_shell tool specification.
- tool-00002-finance-fx-price/: finance_fx_price tool specification.
- tool-00003-skill/: skill tool specification.
- tool-00004-web-fetch/: web_fetch tool specification.
- tool-00005-web-search/: web_search tool specification.
- llm/: Language model adapter specification index.
- template.spec.md: Template for new component or feature specs.

## Delivery Rules

- Every requirement must have a stable ID in the form REQ-<AREA>-<NUMBER>.
- Every acceptance scenario must map to one or more tests.
- New behavior must be represented as scenario updates before code changes.
- Breaking changes require explicit migration notes in the related spec file.
