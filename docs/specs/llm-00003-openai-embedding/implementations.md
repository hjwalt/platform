# Implementation Notes

## Module Layout

```
agent/
├── types.go          # EmbeddingInput, EmbeddingOutput structs; agent.Embedding interface
└── llm/
    ├── types.go      # EmbeddingConfig struct, NewEmbedding factory
    └── embedding.go  # OpenAI embedding adapter (createOpenAiEmbedding, openAiEmbedding struct, Embed method)
```

## Types

### EmbeddingInput

```go
type EmbeddingInput struct {
    Text []string
}
```

Wraps one or more text inputs to be embedded. The adapter maps `Text` to the OpenAI `input` field as an array of strings.

### EmbeddingOutput

```go
type EmbeddingOutput struct {
    Embedding [][]float64
}
```

Wraps one or more embedding vectors. `Embedding[i]` corresponds to `EmbeddingInput.Text[i]`. The adapter populates this from `CreateEmbeddingResponse.Data[].Embedding`.

### Embedding Interface

```go
type Embedding interface {
    runtime.Runtime
    Embed(ctx context.Context, inputs EmbeddingInput) (EmbeddingOutput, error)
}
```

## Adapter Structure

The `openAiEmbedding` struct mirrors the `openAiModel` pattern:

```go
type openAiEmbedding struct {
    Model      string
    Endpoint   string
    Secret     string
    Dimensions int
    client     openai.Client
}
```

- `Start()` initializes `openai.NewClient(option.WithBaseURL(...), option.WithAPIKey(...))`
- `Stop()` is a no-op (matches existing `openAiModel.Stop()`)
- `Embed(ctx, inputs)` maps `inputs.Text` → `EmbeddingNewParamsInputUnion{OfArrayOfStrings: inputs.Text}` and unpacks `CreateEmbeddingResponse.Data[].Embedding` into `EmbeddingOutput.Embedding`

## EmbeddingConfig

```go
type EmbeddingConfig struct {
    Type       ModelType
    Model      string
    Endpoint   string
    Secret     string
    Dimensions int  // 0 = use model default
}
```

## Factory Function

```go
func NewEmbedding(config EmbeddingConfig) agent.Embedding {
    switch config.Type {
    case OpenAi:
        return createOpenAiEmbedding(config)
    default:
        return createOpenAiEmbedding(config)
    }
}
```

## Embed Method — Input/Output Mapping

```go
func (r *openAiEmbedding) Embed(ctx context.Context, inputs agent.EmbeddingInput) (agent.EmbeddingOutput, error) {
    if len(inputs.Text) == 0 {
        return agent.EmbeddingOutput{}, nil
    }

    params := openai.EmbeddingNewParams{
        Input: openai.EmbeddingNewParamsInputUnion{OfArrayOfStrings: inputs.Text},
        Model: openai.EmbeddingModel(r.Model),
    }
    if r.Dimensions > 0 {
        params.Dimensions = param.NewOpt(int64(r.Dimensions))
    }

    resp, err := r.client.Embeddings.New(ctx, params)
    if err != nil {
        return agent.EmbeddingOutput{}, err
    }

    embeddings := make([][]float64, len(resp.Data))
    for i, e := range resp.Data {
        embeddings[i] = e.Embedding
    }
    return agent.EmbeddingOutput{Embedding: embeddings}, nil
}
```

## Dependencies

- `github.com/openai/openai-go/v3` (already in go.mod)
- `github.com/hjwalt/platform/agent` (interface definition)
- `github.com/hjwalt/platform/runtime` (Runtime interface)

## Edge Cases

| Case | Behavior |
|------|----------|
| Empty Text slice | Return zero-valued EmbeddingOutput, no API call |
| Single text | Single-element batch → single-element result |
| Large batch (>2048 items) | OpenAI API rejects; error returned to caller |
| Context cancelled | Embed returns ctx.Err() |
| API returns fewer results than inputs | Return mismatch error |
| Empty string in Text | OpenAI API rejects; error returned |
