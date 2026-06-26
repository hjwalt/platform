package llm

import (
	"context"

	"github.com/hjwalt/platform/agent"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
)

func createOpenAiEmbedding(config EmbeddingConfig) agent.Embedding {
	return &openAiEmbedding{
		Model:      config.Model,
		Endpoint:   config.Endpoint,
		Secret:     config.Secret,
		Dimensions: config.Dimensions,
	}
}

type openAiEmbedding struct {
	Model      string
	Endpoint   string
	Secret     string
	Dimensions int
	client     openai.Client
}

func (r *openAiEmbedding) Start() error {
	r.client = openai.NewClient(
		option.WithBaseURL(r.Endpoint),
		option.WithAPIKey(r.Secret),
	)

	return nil
}

func (r *openAiEmbedding) Stop() {
}

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
