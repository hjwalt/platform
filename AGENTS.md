# AGENTS.md

## Repository Structure

This is a Go monorepo with the following main packages:

- `agent` - LLM/tool harness and agent runtime components
- `flow` - Current dataflow primitives and runtime integrations
- `flows` - Legacy dataflow package (reference only for migration context)
- `format` - Serialization and byte masking/encryption helpers
- `reflect` - Safe dynamic type conversion helpers
- `runtime` - Lifecycle management for long-running services
- `state` - Key-value state interfaces and implementations
- `web` - HTTP/web component utilities
- `main.go` - Root entrypoint

## Key Commands

### Development

- `make test` - Run tests
- `make tidy` - Clean up go.mod
- `make update` - Update dependencies
- `make mocks` - Generate mocks
- `make proto` - Generate protobuf files
- `make cov` - Run coverage
- `make htmlcov` - Generate HTML coverage report

### Build and Runtime

- `make build` - Build binary to `bin/platform`
- `make run` - Run example application via script
- `make reset` - Reset topics/tables and seed demo data

### Local Services

- `make up` - Start local dependencies with Docker Compose
- `make down` - Stop local dependencies

### Running Examples

- `docker compose up -d` - Start required services (Kafka, Postgres)
- `make reset` - Clean up topics and tables
- `make run` - Start example application

## Progressive Disclosure

Read only what is needed for the current task from the following:

- [Architecture](docs/memory/architecture.md)
- [Spec Driven Development](docs/memory/spec-driven-development.md)

## Testing

Tests use `testcontainers-go` which requires rootless Podman setup. Configure with:

```
export DOCKER_HOST=unix:///run/user/1000/podman/podman.sock
```

## Special Considerations

- The main entrypoint is in the `main.go` file at the root
- Examples are in the `example` directory
- Container-based services are required for integration tests
- Environment variables are loaded through the `environment` package
- When finishing up any task, check if there are documents in `/docs` that needs updating and update as necessary
