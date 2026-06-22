# Spec Driven Development (SDD)

## How To Use

1. Pick a target SDD spec folder.
2. Select one requirement ID not yet implemented.
3. Write or update failing tests first.
4. Implement only enough code to make tests pass.
5. Run make test and update the requirement status.

## Spec Status Legend

- Proposed: drafted but not started.
- In Progress: tests and implementation are actively being developed.
- Implemented: acceptance tests pass.
- Verified: implemented and explicitly covered in regression tests.

## Rules

1. Each specification must be created in its own folder
2. The folder must be named in this pattern: `<type>-<index>-<short title>`
3. Type can be one of the following:
   - agent
   - backend
   - flow
   - llm
   - repository
   - tool
   - web-component
   - web-page
4. Index is five digit zero padded starting with 1 for every type
5. Short title should include 1 to 5 words to help developers quickly know what the specification is about
6. File templates are in `/docs/templates`
7. Always update `tasks.md`, `ammendments.md` and `implementations.md` during implementation stage
8. Every acceptance scenario must map to one or more tests.
9. New behavior must be represented as scenario updates before code changes.
10. Breaking changes require explicit migration notes in the related spec file.
11. Update the spec index in `/docs/memory/architecture` `spec index` section after a new spec is created or if existing spec is renamed

## Files

All specs reside in `/docs/specs` directory. Each specification should contain these files:

### specs.md

This file contains the specifications for the requirements. It must contain the following sections:

1. Title
2. High Level Description
3. User Scenarios
4. Functional Requirements
5. Non-Functional Requirements
6. Definition of Done
7. Testing Methodology

### tasks.md

This file contains the to-do list for the agents to complete to fully develop for the specifications. It must contain the following sections:

1. Preparation
2. Implementation
3. Validation

Each task should follow the following format:

- [ ] task description

After tasks are performed, fill the [ ] with an x like [x].

### ammendments.md

This file contains the specfication ammendment history with numbered sequence

### implementations.md

This file contains the implemenation details for the feature. It must contain the following sections:

1. Choices Made
2. Libraries Used
3. Implementation Preferences
4. Caveats

## Spec Index

- agent-00001-harness-contracts/: SDD template-based agent harness contracts specification.
- backend-00001-state-file/: File State Store specification.
- backend-00002-state-memory/: In-Memory State Store specification.
- flow-00001-flow-runtime/: SDD template-based flow/runtime lifecycle specification.
- llm-00001-deepseek/: DeepSeek LLM adapter specification.
- llm-00002-openai/: OpenAI LLM adapter specification.
- llm-00003-openai-embedding/: OpenAI Embedding adapter specification.
- repository-00001-repository-foundation/: Repository-wide engineering constraints and delivery gates.
- tool-00001-linux-shell/: linux_shell tool specification.
- tool-00002-finance-fx-price/: finance_fx_price tool specification.
- tool-00003-subagent/: subagent tool specification.
- tool-00004-web-fetch/: web_fetch tool specification.
- tool-00005-web-search/: web_search tool specification.
- tool-00006-finance-stock-price/: finance_stock_price tool specification.
- tool-00007-memory-management/: unified memory tool specification with operation=get|update|clear.
- tool-00008-tool-search/: tool_search tool specification — natural language search over registered tools.
- web-component-00001-chat-item/: Chat Item web component specification.
- web-component-00002-chat-list/: Chat List web component specification.
- web-component-00003-sidebar/: Sidebar web component specification.
- web-component-00004-sidebar-button/: Sidebar Button web component specification.
- web-component-00005-sidebar-button-list/: Sidebar Button List web component specification.
- web-component-00006-sidebar-item/: Sidebar Item web component specification.
- web-component-00007-sidebar-item-header/: Sidebar Item Header web component specification.
- web-component-00008-sidebar-item-list/: Sidebar Item List web component specification.
- web-page-00001-home/: page_home route and dashboard rendering specification.
- web-page-00002-chat/: page_chat routing, posting, and tool-decision specification.
- web-page-00003-billing/: page_billing route and dashboard rendering specification.
- web-page-00004-error-500/: page_error_500 fallback rendering specification.
