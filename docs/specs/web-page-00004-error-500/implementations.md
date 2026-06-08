# Choices Made

- Keep a standalone Error function as the reusable fallback entrypoint.
- Preserve embedded template rendering through Html and render.Component.

# Libraries Used

- Go standard library packages embed and net/http.
- Internal package github.com/hjwalt/platform/web/render.

# Implementation Preferences

- Keep fallback rendering independent from page-specific models.
- Avoid side effects in error rendering paths.

# Caveats

- Current implementation does not expose dynamic error details in the view.
- Template-only failures may surface at runtime if embedding or parsing changes.
