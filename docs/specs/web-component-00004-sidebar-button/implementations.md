# Choices Made

- Preserve existing package-level View API for web/component/component_sidebar_button.
- Keep template embedding via go:embed and component.html.

# Libraries Used

- Go standard library packages embed and html/template.
- Internal package github.com/hjwalt/platform/web/render.

# Implementation Preferences

- Prefer explicit model fields over map-based payloads.
- Keep component composition shallow and predictable.

# Caveats

- Rendered HTML assertions may be brittle if templates are heavily reformatted.
- Composition contracts must stay synchronized with parent components.
