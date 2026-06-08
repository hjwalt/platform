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

## Architecture Notes

### Package Boundaries

- `example` package contains runnable usage examples
- `flow` implements dataflow processing
- `flows` is a legacy package for older dataflow processing; avoid adding new code there
- `format` contains data formats for serialization/deserialization and masking/encryption
- `reflect` handles conversion of data types dynamically and safely
- `web` handles HTTP routing and web components
- `state` provides key value pair based state management for stateful functions

### Runtime Management

The platform uses a runtime system for managing long-running services like HTTP servers and message queue consumers. Runtimes are started with `runtime.Start()` and waited for with `runtime.Wait()`.

### Dataflow Processing

Dataflow capabilities are centered in `flow` (current) and `flows` (legacy reference) and follow a Kappa Architecture style with:

- Stateless functions for basic operations
- Stateful functions for aggregations
- Join patterns using intermediate topics
- Materialisation for database upserts
- Task processing for long-running operations

### Environment Configuration

Uses `github.com/hjwalt/platform/environment` for environment variable handling with default values.

## Testing

Tests use `testcontainers-go` which requires rootless Podman setup. Configure with:

```
export DOCKER_HOST=unix:///run/user/1000/podman/podman.sock
```

## Key Components

- State management: `github.com/hjwalt/platform/state`
- Flow execution: `github.com/hjwalt/platform/flow`

## Special Considerations

- The main entrypoint is in the `main.go` file at the root
- Examples are in the `example` directory
- Container-based services are required for integration tests
- Environment variables are loaded through the `environment` package

## Spec Driven Development

### Rules

1. Each specification must be created in its own folder
2. The folder must be named in this pattern: `<type>-<index>-<short title>`
3. Type can be one of the following: `agent`, `web-page`, `web-component`, `backend`, `tool`, `llm`, `flow`. Suggest new type if the current set is not accurate enough
4. Index is five digit zero padded starting with 1 for every type
5. Short title should include 1 to 5 words to help developers quickly know what the specification is about
6. File templates are in `docs/templates`

### Files

Each specification should contain these files:

#### specs.md

This file contains the specifications for the requirements. It must contain the following sections:

1. Title
2. High Level Description
3. User Scenarios
4. Functional Requirements
5. Non-Functional Requirements
6. Definition of Done
7. Testing Methodology

#### tasks.md

This file contains the to-do list for the agents to complete to fully develop for the specifications. It must contain the following sections:

1. Preparation
2. Implementation
3. Validation

Each task should follow the following format:

- [ ] task description

After tasks are performed, fill the [ ] with an x like [x].

#### ammendments.md

This file contains the specfication ammendment history with numbered sequence

#### implementations.md

This file contains the implemenation details for the feature. It must contain the following sections:

1. Choices Made
2. Libraries Used
3. Implementation Preferences
4. Caveats
