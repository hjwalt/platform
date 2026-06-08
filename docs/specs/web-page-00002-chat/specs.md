# Title

Web Page Spec: Chat

# High Level Description

Define the behavior, routing, interaction flow, and rendering contract for web/page/page_chat, including chat context selection, message posting, and tool approval or rejection actions.

# User Scenarios

1. As a user, I can open the chat page and switch among available chat contexts.
2. As a user, I can submit a chat message and see the chat view update.
3. As a user, I can approve or reject tool execution requests from the chat UI.
4. As a maintainer, I can extend chat page behavior while keeping route and message contracts stable.

## Functional Requirements

### FR-CHAT-001 Route Registration

1. The package registers GET /chat and GET /chat/{chat_id} routes for page rendering.
2. The package registers POST /chat/{chat_id} for user message submission.
3. The package registers POST /chat/{chat_id}/accept, POST /chat/{chat_id}/reject, and POST /chat-view for tool and context actions.
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

# Non-Functional Requirements

1. Route behavior is deterministic for identical request payloads.
2. Handler logic remains isolated to web/page/page_chat without global mutable state.
3. Chat rendering and post handlers remain compatible with HTMX form interaction patterns in page.html.

# Definition of Done

1. This spec accurately reflects route and handler behavior in page_chat.
2. Tests or validation steps cover route wiring and key post flows.
3. Message publication contracts for user and tool actions are validated.

# Testing Methodology

1. Add or update tests for GET and POST route behavior, including chat_id selection.
2. Validate form parsing failure behavior for post handlers.
3. Run make test and perform a manual chat flow smoke test for submit, accept, and reject actions.
