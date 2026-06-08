# Tasks

## Preparation

- [ ] Review route constants and handler flow in web/page/page_chat.
- [ ] Confirm expected message and tool action payload fields from form submissions.
- [ ] Confirm reset interaction design in page.html and expected store clearing behavior.

## Implementation

- [ ] Align Add route registration with FR-CHAT-001 requirements.
- [ ] Add or update tests for get, getWithId, post, postTool, rejectTool, and postChatView.
- [ ] Verify rendering composition for chat list and sidebar layout contract.
- [ ] Add reset button UI wiring for the active chat context.
- [ ] Implement reset handler to clear chat history and return updated chat list view.
- [ ] Add or update tests for reset route and reset handler behavior.

## Validation

- [ ] Run make test and confirm affected package tests pass.
- [ ] Perform manual chat flow checks for context switch and message submit.
- [ ] Validate tool accept and reject flows produce expected message types.
- [ ] Validate reset button flow clears history and remains stable on repeated reset.
