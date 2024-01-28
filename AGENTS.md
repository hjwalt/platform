# AGENTS.md

## Repository Structure

This is a Go monorepo with the following main packages:
- `commons` - Core utilities and infrastructure
- `routes` - HTTP routing and web components
- `flows` - Dataflow processing with Kafka and RabbitMQ
- `main` - Entry point for the web application

## Key Commands

### Development
- `make test` - Run tests
- `make tidy` - Clean up go.mod
- `make update` - Update dependencies
- `make mocks` - Generate mocks
- `make proto` - Generate protobuf files
- `make cov` - Run coverage
- `make htmlcov` - Generate HTML coverage report

### Running Examples
- `docker-compose up -d` - Start required services (Kafka, Postgres)
- `make reset` - Clean up topics and tables
- `make run` - Start example application

## Architecture Notes

### Package Boundaries
- `commons` is the core library package
- `examples` package for adding examples on how to use the repository
- `flow` implements dataflow processing
- `flows` legacy package for dataflow processing, it can be used as a reference but no new code should be updated here
- `format` contains the dataformat for serialisation, deserialisation, and encryption or byte masking
- `reflect` handles conversion of data types dynamically and safely
- `routes` handles HTTP routing and web components
- `state` provides key value pair based state management for stateful functions

### Runtime Management
The platform uses a runtime system for managing long-running services like HTTP servers and message queue consumers. Runtimes are started with `runtime.Start()` and waited for with `runtime.Wait()`.

### Dataflow Processing
The `flows` package implements a Kappa Architecture approach with:
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
- Examples are in the `examples` directory
- Container-based services are required for integration tests
- Environment variables are loaded through the `commons/environment` package
- Ignore http related components until we completely refactor the `flows` package to `flow` package