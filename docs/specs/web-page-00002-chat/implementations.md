# Choices Made

- Preserve explicit route constants for chat, id, accept, reject, and chat-view paths.
- Keep message and tool workflows coupled to AgentMessageProducer contracts.
- Add a reset action per chat context using a dedicated reset route and form action.

# Libraries Used

- Go standard library packages embed, net/http, and log/slog.
- External router package github.com/go-chi/chi/v5.
- Internal packages agent, web/layout, web/render, and chat components.

# Implementation Preferences

- Prefer explicit validation of required form fields before publishing messages.
- Keep handler return views minimal for HTMX-driven partial updates.
- Keep reset behavior idempotent and scoped to the active chat context only.

# Caveats

- Current handlers ignore store read/key errors and may mask data-loading failures.
- Chat behavior depends on downstream AgentHarnessStore and AgentMessageProducer reliability.
- Reset behavior requires store-level clear/delete semantics for a single chat key.
