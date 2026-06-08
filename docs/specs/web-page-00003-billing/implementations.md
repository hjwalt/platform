# Choices Made

- Preserve billing route registration via Add and GET /billing.
- Keep composition through layout.Dashboard with sidebar and embedded template content.

# Libraries Used

- Go standard library packages embed and net/http.
- Internal packages web/layout, web/render, and web/route.

# Implementation Preferences

- Prefer explicit model struct definition, even when no current fields are required.
- Keep render.Component arguments deterministic and straightforward.

# Caveats

- Billing template assertions may be sensitive to markup-only updates.
- Shared dashboard layout changes can affect billing page output.
