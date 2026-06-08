# Title

Web Page Spec: Chat

# High Level Description

Define the behavior, routing, interaction flow, and rendering contract for web/page/page_chat, including chat context selection, message posting, tool approval or rejection actions, and chat history reset from the UI.

# User Scenarios

1. As a user, I can open the chat page and switch among available chat contexts.
2. As a user, I can submit a chat message and see the chat view update.
3. As a user, I can approve or reject tool execution requests from the chat UI.
4. As a user, I can reset the current chat history using a dedicated reset button.
5. As a maintainer, I can extend chat page behavior while keeping route and message contracts stable.

## Functional Requirements

### FR-CHAT-001 Route Registration

1. The package registers GET /chat and GET /chat/{chat_id} routes for page rendering.
2. The package registers POST /chat/{chat_id} for user message submission.
3. The package registers POST /chat/{chat_id}/accept, POST /chat/{chat_id}/reject, POST /chat/{chat_id}/reset, and POST /chat-view for tool and context actions.
4. All routes are wrapped with render.Handle and page_error_500.Error.

### FR-CHAT-002 Page Rendering and Context Loading

1. The default GET handler resolves to chat context web.
2. Chat-specific GET handler resolves chat_id from the route path.
3. View rendering uses layout.Dashboard and component_sidebar.View.
4. The chat list region is rendered through component_chat_list and chat item components.

### FR-CHAT-003 Message Submission Behavior

1. POST chat handlers parse form data and return HTTP 400 on parse failures.
2. When message is present, the package publishes a user message through AgentMessageProducer.
3. When message is missing, the package returns a placeholder chat item indicating no message received.

### FR-CHAT-004 Tool Decision Behavior

1. Tool accept action publishes MessageType_ToolExecute with the selected tool call payload.
2. Tool reject action publishes MessageType_ToolResult with rejection details.
3. Missing required tool fields produce a fallback no message received chat item.

### FR-CHAT-005 Reset Chat History Behavior

1. The chat page UI exposes a reset button for the active chat context.
2. Reset action submits a POST request to /chat/{chat_id}/reset for the active context.
3. Reset action clears stored message history for the selected chat context.
4. After reset, the chat list view returns an empty history state for that context.
5. Reset failures return an error response consistent with existing form/handler patterns.

# Non-Functional Requirements

1. Route behavior is deterministic for identical request payloads.
2. Handler logic remains isolated to web/page/page_chat without global mutable state.
3. Chat rendering and post handlers remain compatible with HTMX form interaction patterns in page.html.
4. Reset interactions are idempotent so repeated resets do not produce inconsistent state.

# Definition of Done

1. This spec accurately reflects route and handler behavior in page_chat.
2. Tests or validation steps cover route wiring and key post flows.
3. Message publication contracts for user and tool actions are validated.
4. Reset button visibility, route handling, and history-clearing behavior are validated.

# Testing Methodology

1. Add or update tests for GET and POST route behavior, including chat_id selection.
2. Validate form parsing failure behavior for post handlers.
3. Validate reset handler behavior for existing and already-empty chat contexts.
4. Run make test and perform a manual chat flow smoke test for submit, accept, reject, and reset actions.
