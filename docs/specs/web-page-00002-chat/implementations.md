# Choices Made

- Preserve explicit route constants for chat, id, accept, reject, and chat-view paths.
- Keep message and tool workflows coupled to AgentMessageProducer contracts.

# Libraries Used

- Go standard library packages embed, net/http, and log/slog.
- External router package github.com/go-chi/chi/v5.
- Internal packages agent, web/layout, web/render, and chat components.

# Implementation Preferences

- Prefer explicit validation of required form fields before publishing messages.
- Keep handler return views minimal for HTMX-driven partial updates.

# Caveats

- Current handlers ignore store read/key errors and may mask data-loading failures.
- Chat behavior depends on downstream AgentHarnessStore and AgentMessageProducer reliability.
