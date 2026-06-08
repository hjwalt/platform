# Tasks

## Preparation

- [x] Review route constants and handler flow in web/page/page_chat.
- [x] Confirm expected message and tool action payload fields from form submissions.
- [x] Confirm reset interaction design in page.html and expected store clearing behavior.

## Implementation

- [x] Align Add route registration with FR-CHAT-001 requirements.
- [ ] Add or update tests for get, getWithId, post, postTool, rejectTool, and postChatView.
- [x] Verify rendering composition for chat list and sidebar layout contract.
- [x] Add reset button UI wiring for the active chat context.
- [x] Implement reset handler to clear chat history and return updated chat list view.
- [x] Add or update tests for reset route and reset handler behavior.

## Validation

- [x] Run make test and confirm affected package tests pass.
- [ ] Perform manual chat flow checks for context switch and message submit.
- [ ] Validate tool accept and reject flows produce expected message types.
- [ ] Validate reset button flow clears history and remains stable on repeated reset.
