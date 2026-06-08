# Repository Specs For Spec-Driven Development

This directory defines executable, test-first specifications for the platform repository.

## How To Use

1. Pick a target spec file or SDD spec folder.
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

- repository-00001-repository-foundation/: Repository-wide engineering constraints and delivery gates.
- flow-00001-flow-runtime/: SDD template-based flow/runtime lifecycle specification.
- agent-00001-harness-contracts/: SDD template-based agent harness contracts specification.
- tool-00001-linux-shell/: linux_shell tool specification.
- tool-00002-finance-fx-price/: finance_fx_price tool specification.
- tool-00003-skill/: skill tool specification.
- tool-00004-web-fetch/: web_fetch tool specification.
- tool-00005-web-search/: web_search tool specification.
- tool-00006-finance-stock-price/: finance_stock_price tool specification.
- tool-00007-memory-management/: memory_get, memory_update, and memory_clear tool specification.
- web-page-00001-home/: page_home route and dashboard rendering specification.
- web-page-00002-chat/: page_chat routing, posting, and tool-decision specification.
- web-page-00003-billing/: page_billing route and dashboard rendering specification.
- web-page-00004-error-500/: page_error_500 fallback rendering specification.

## Delivery Rules

- Every requirement must have a stable ID in the form REQ-<AREA>-<NUMBER>.
- Every acceptance scenario must map to one or more tests.
- New behavior must be represented as scenario updates before code changes.
- Breaking changes require explicit migration notes in the related spec file.
