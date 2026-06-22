# Amendments

## 2026-06-22 — Rename to openai-embedding, add input/output structs

- Renamed spec from `llm-00003-embedding` to `llm-00003-openai-embedding` for consistency with `llm-00001-deepseek` and `llm-00002-openai` naming.
- Added `agent.EmbeddingInput` struct with `Text []string` field, bundling batch inputs into a single struct.
- Added `agent.EmbeddingOutput` struct with `Embedding [][]float64` field, bundling batch results into a single struct.
- Defined `Embedding` interface as `Embed(ctx context.Context, inputs EmbeddingInput) (EmbeddingOutput, error)` — accepting and returning single structs rather than slices.

## Initial Spec — 2026-06-22

- Created specification for OpenAI embedding adapter.
- Added `agent.Embedding` interface to `agent/types.go`.
- Defined `EmbeddingConfig`, `NewEmbedding` factory, and `createOpenAiEmbedding` constructor patterns.
- Seven functional requirements covering lifecycle, single/batch embed, model/dimension config, error handling, and interface compliance.
