# Title

OpenAI Embedding Adapter Specification

# High Level Description

This specification defines expected behavior for the OpenAI embedding adapter in agent/llm/embedding.go.
The adapter must implement the shared `agent.Embedding` lifecycle and embed contract, accepting `agent.EmbeddingInput` and returning `agent.EmbeddingOutput` via the OpenAI embeddings API.

# User Scenarios

1. As an agent runtime, I want to initialize the OpenAI embedding client from config so vectors can be generated against a configured endpoint.
2. As a memory or retrieval system, I want to generate embedding vectors for one or more texts via `EmbeddingInput` so they can be stored and compared.
3. As a batch processing pipeline, I want to generate vectors for multiple texts in a single API call so throughput improves.
4. As an operator, I want transport or API failures surfaced as an error return so callers can retry or degrade.

# Functional Requirements

## FR-EMB-001 Model Construction And Startup

- createOpenAiEmbedding must map EmbeddingConfig fields into the adapter instance.
- Start must create an OpenAI client using endpoint and secret from config.
- Stop must be safe to call and must not panic.

## FR-EMB-002 Model Selection

- Embed must use the model name from EmbeddingConfig when constructing the API request.
- The adapter must support model override through configuration without code changes.

## FR-EMB-003 Dimension Configuration

- When EmbeddingConfig.Dimensions is set, Embed must pass it to the API request.
- When Dimensions is zero, Embed must omit the parameter and rely on the model default.

## FR-EMB-004 Error Handling

- If the embedding API request fails, Embed must return a non-nil error.
- On failure, the returned `EmbeddingOutput` must be zero-valued.
- The error must wrap the underlying API or transport error for upstream handling.

## FR-EMB-005 Interface Compliance

- The adapter must implement `agent.Embedding` (extending `runtime.Runtime`).
- Embed must accept `context.Context` and honor context cancellation or deadline.
- The adapter must remain safe for concurrent use after Start completes.

# Non-Functional Requirements

1. Reliability: Embed must tolerate empty inputs without panic and return empty results cleanly.
2. Determinism: Identical inputs against the same model must produce identical vectors (API-dependent).
3. Maintainability: The adapter constructor and mapping logic must remain independently unit-testable.
4. Observability: API-level errors must be logged or wrapped in a way that supports upstream diagnosis.

# Definition of Done

1. FR-EMB-001 through FR-EMB-007 are covered by automated tests.
2. `make test` passes.
3. `tasks.md` is fully checked.
4. Any behavior changes are recorded in repository amendment/changelog process.

# Testing Methodology

1. Unit tests in agent/llm for startup, single-embed, batch-embed, and error-path behavior.
2. Use mocked or stubbed OpenAI embedding transport/client behavior for success and failure paths.
3. Verify dimension and model parameters are forwarded correctly through API call assertions.
4. Validate empty-input handling and result ordering for batch calls.
5. Test context cancellation propagation.
