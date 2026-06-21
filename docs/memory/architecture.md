# Architecture Notes

## Package Boundaries

- `example` package contains runnable usage examples
- `flow` implements dataflow processing
- `format` contains data formats for serialization/deserialization and masking/encryption
- `reflect` handles conversion of data types dynamically and safely
- `web` handles HTTP routing and web components
- `state` provides key value pair based state management for stateful functions

## Runtime Management

The platform uses a runtime system for managing long-running services like HTTP servers and message queue consumers. Runtimes are started with `runtime.Start()` and waited for with `runtime.Wait()`.

## Dataflow Processing

Dataflow capabilities are centered in `flow` (current) and `flows` (legacy reference) and follow a Kappa Architecture style with:

- Stateless functions for basic operations
- Stateful functions for aggregations
- Join patterns using intermediate topics
- Materialisation for database upserts
- Task processing for long-running operations

## Environment Configuration

Uses `github.com/hjwalt/platform/environment` for environment variable handling with default values.
