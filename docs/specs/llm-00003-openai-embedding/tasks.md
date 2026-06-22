# Tasks

## Preparation

- [ ] Confirm embedding adapter requirement IDs and wording with maintainers.
- [ ] Review OpenAI embeddings API v3 SDK surface for compatibility with proposed interface.

## Implementation

- [ ] Add `EmbeddingInput` struct (`Text []string`) and `EmbeddingOutput` struct (`Embedding [][]float64`) to `agent/types.go`.
- [ ] Add `Embedding` interface to `agent/types.go` extending `runtime.Runtime` with `Embed(ctx, EmbeddingInput) (EmbeddingOutput, error)`.
- [ ] Add `EmbeddingConfig` struct to `agent/llm/types.go` with Type, Model, Endpoint, Secret, and Dimensions fields.
- [ ] Add `createOpenAiEmbedding` constructor in `agent/llm/embedding.go`.
- [ ] Add startup and lifecycle tests for createOpenAiEmbedding/Start/Stop (FR-EMB-001).
- [ ] Implement and test single-input embedding (FR-EMB-002).
- [ ] Implement and test batch-input embedding (FR-EMB-003).
- [ ] Implement and test model selection forwarding (FR-EMB-004).
- [ ] Implement and test dimension configuration forwarding (FR-EMB-005).
- [ ] Implement and test error-path behavior for API and transport failures (FR-EMB-006).
- [ ] Add interface compliance tests verifying runtime.Runtime and agent.Embedding satisfaction (FR-EMB-007).
- [ ] Add `NewEmbedding` factory function to `agent/llm/types.go` wiring config type to constructor.

## Validation

- [ ] Run `make test` and fix failures.
- [ ] Update amendment/changelog artifacts per repo process.
- [ ] Mark completed tasks and document residual risks.
