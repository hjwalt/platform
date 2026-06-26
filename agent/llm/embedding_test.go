package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockEmbeddingServer returns an httptest server that responds with the given
// embedding vectors. Each []float64 in embeddings becomes one data entry.
func mockEmbeddingServer(t *testing.T, embeddings [][]float64, statusCode int) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "mock error",
			})
			return
		}

		data := make([]map[string]interface{}, len(embeddings))
		for i, emb := range embeddings {
			data[i] = map[string]interface{}{
				"object":    "embedding",
				"index":     i,
				"embedding": emb,
			}
		}

		resp := map[string]interface{}{
			"object": "list",
			"data":   data,
			"model":  "text-embedding-3-small",
			"usage": map[string]interface{}{
				"prompt_tokens": 8,
				"total_tokens":  8,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))

	t.Cleanup(server.Close)
	return server
}

func newTestEmbedding(serverURL string) *openAiEmbedding {
	return &openAiEmbedding{
		Model:      "text-embedding-3-small",
		Endpoint:   serverURL,
		Secret:     "test-secret",
		Dimensions: 0,
	}
}

// FR-EMB-007: Interface compliance
func TestOpenAiEmbeddingSatisfiesInterface(t *testing.T) {
	var emb agent.Embedding = &openAiEmbedding{}
	assert.NotNil(t, emb)

	var rt runtime.Runtime = &openAiEmbedding{}
	assert.NotNil(t, rt)
}

// FR-EMB-001: Model Construction And Startup
func TestCreateOpenAiEmbedding(t *testing.T) {
	config := EmbeddingConfig{
		Type:       OpenAi,
		Model:      "text-embedding-3-small",
		Endpoint:   "https://api.openai.com/v1",
		Secret:     "sk-test",
		Dimensions: 1536,
	}

	emb := NewEmbedding(config)
	require.NotNil(t, emb)

	openAiEmb, ok := emb.(*openAiEmbedding)
	require.True(t, ok)
	assert.Equal(t, "text-embedding-3-small", openAiEmb.Model)
	assert.Equal(t, "https://api.openai.com/v1", openAiEmb.Endpoint)
	assert.Equal(t, "sk-test", openAiEmb.Secret)
	assert.Equal(t, 1536, openAiEmb.Dimensions)
}

func TestStartAndStop(t *testing.T) {
	server := mockEmbeddingServer(t, nil, http.StatusOK)
	emb := newTestEmbedding(server.URL)

	err := emb.Start()
	assert.NoError(t, err)

	// Stop should be safe and not panic
	emb.Stop()
}

func TestStopIsIdempotent(t *testing.T) {
	server := mockEmbeddingServer(t, nil, http.StatusOK)
	emb := newTestEmbedding(server.URL)

	require.NoError(t, emb.Start())
	emb.Stop()
	emb.Stop() // should not panic
	emb.Stop()
}

// FR-EMB-002 / FR-EMB-003: Single and batch embedding
func TestEmbedSingleText(t *testing.T) {
	expected := []float64{0.1, 0.2, 0.3}
	server := mockEmbeddingServer(t, [][]float64{expected}, http.StatusOK)
	emb := newTestEmbedding(server.URL)
	require.NoError(t, emb.Start())

	ctx := context.Background()
	output, err := emb.Embed(ctx, agent.EmbeddingInput{Text: []string{"hello"}})
	require.NoError(t, err)
	assert.Len(t, output.Embedding, 1)
	assert.Equal(t, expected, output.Embedding[0])
}

func TestEmbedBatchText(t *testing.T) {
	expected := [][]float64{
		{0.1, 0.2, 0.3},
		{0.4, 0.5, 0.6},
		{0.7, 0.8, 0.9},
	}
	server := mockEmbeddingServer(t, expected, http.StatusOK)
	emb := newTestEmbedding(server.URL)
	require.NoError(t, emb.Start())

	ctx := context.Background()
	output, err := emb.Embed(ctx, agent.EmbeddingInput{
		Text: []string{"hello", "world", "foo"},
	})
	require.NoError(t, err)
	assert.Len(t, output.Embedding, 3)
	assert.Equal(t, expected[0], output.Embedding[0])
	assert.Equal(t, expected[1], output.Embedding[1])
	assert.Equal(t, expected[2], output.Embedding[2])
}

// FR-EMB-003: Empty input handling
func TestEmbedEmptyInput(t *testing.T) {
	server := mockEmbeddingServer(t, nil, http.StatusOK)
	emb := newTestEmbedding(server.URL)
	require.NoError(t, emb.Start())

	ctx := context.Background()
	output, err := emb.Embed(ctx, agent.EmbeddingInput{Text: []string{}})
	require.NoError(t, err)
	assert.Empty(t, output.Embedding)
}

func TestEmbedNilInput(t *testing.T) {
	server := mockEmbeddingServer(t, nil, http.StatusOK)
	emb := newTestEmbedding(server.URL)
	require.NoError(t, emb.Start())

	ctx := context.Background()
	output, err := emb.Embed(ctx, agent.EmbeddingInput{Text: nil})
	require.NoError(t, err)
	assert.Empty(t, output.Embedding)
}

// FR-EMB-004: Model selection forwarding
func TestEmbedModelForwarding(t *testing.T) {
	var capturedModel string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if m, ok := body["model"]; ok {
			capturedModel = m.(string)
		}

		resp := map[string]interface{}{
			"object": "list",
			"data": []map[string]interface{}{
				{"object": "embedding", "index": 0, "embedding": []float64{0.1}},
			},
			"model":  capturedModel,
			"usage":  map[string]interface{}{"prompt_tokens": 1, "total_tokens": 1},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(server.Close)

	emb := &openAiEmbedding{
		Model:    "text-embedding-ada-002",
		Endpoint: server.URL,
		Secret:   "test-secret",
	}
	require.NoError(t, emb.Start())

	ctx := context.Background()
	_, err := emb.Embed(ctx, agent.EmbeddingInput{Text: []string{"test"}})
	require.NoError(t, err)
	assert.Equal(t, "text-embedding-ada-002", capturedModel)
}

// FR-EMB-005: Dimension configuration forwarding
func TestEmbedDimensionsForwarding(t *testing.T) {
	var capturedDimensions float64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if d, ok := body["dimensions"]; ok {
			capturedDimensions = d.(float64)
		}

		resp := map[string]interface{}{
			"object": "list",
			"data": []map[string]interface{}{
				{"object": "embedding", "index": 0, "embedding": make([]float64, 256)},
			},
			"model":  "text-embedding-3-small",
			"usage":  map[string]interface{}{"prompt_tokens": 1, "total_tokens": 1},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(server.Close)

	emb := &openAiEmbedding{
		Model:      "text-embedding-3-small",
		Endpoint:   server.URL,
		Secret:     "test-secret",
		Dimensions: 256,
	}
	require.NoError(t, emb.Start())

	ctx := context.Background()
	_, err := emb.Embed(ctx, agent.EmbeddingInput{Text: []string{"test"}})
	require.NoError(t, err)
	assert.Equal(t, float64(256), capturedDimensions)
}

func TestEmbedDimensionsOmittedWhenZero(t *testing.T) {
	var hasDimensions bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		_, hasDimensions = body["dimensions"]

		resp := map[string]interface{}{
			"object": "list",
			"data": []map[string]interface{}{
				{"object": "embedding", "index": 0, "embedding": []float64{0.1}},
			},
			"model":  "text-embedding-3-small",
			"usage":  map[string]interface{}{"prompt_tokens": 1, "total_tokens": 1},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(server.Close)

	emb := &openAiEmbedding{
		Model:      "text-embedding-3-small",
		Endpoint:   server.URL,
		Secret:     "test-secret",
		Dimensions: 0, // zero → omit
	}
	require.NoError(t, emb.Start())

	ctx := context.Background()
	_, err := emb.Embed(ctx, agent.EmbeddingInput{Text: []string{"test"}})
	require.NoError(t, err)
	assert.False(t, hasDimensions, "dimensions should be omitted when zero")
}

// FR-EMB-006: Error handling
func TestEmbedAPIError(t *testing.T) {
	server := mockEmbeddingServer(t, nil, http.StatusInternalServerError)
	emb := newTestEmbedding(server.URL)
	require.NoError(t, emb.Start())

	ctx := context.Background()
	output, err := emb.Embed(ctx, agent.EmbeddingInput{Text: []string{"test"}})
	assert.Error(t, err)
	assert.Empty(t, output.Embedding)
}

func TestEmbedBadGatewayError(t *testing.T) {
	server := mockEmbeddingServer(t, nil, http.StatusBadGateway)
	emb := newTestEmbedding(server.URL)
	require.NoError(t, emb.Start())

	ctx := context.Background()
	output, err := emb.Embed(ctx, agent.EmbeddingInput{Text: []string{"test"}})
	assert.Error(t, err)
	assert.Empty(t, output.Embedding)
	// On error, output should be zero-valued
	assert.Nil(t, output.Embedding)
}

func TestEmbedTransportError(t *testing.T) {
	// Use a non-routable address to trigger a transport error
	emb := &openAiEmbedding{
		Model:    "text-embedding-3-small",
		Endpoint: "http://127.0.0.1:1", // no server listening
		Secret:   "test-secret",
	}
	require.NoError(t, emb.Start())

	ctx := context.Background()
	output, err := emb.Embed(ctx, agent.EmbeddingInput{Text: []string{"test"}})
	assert.Error(t, err)
	assert.Empty(t, output.Embedding)
}

// Context cancellation
func TestEmbedContextCancellation(t *testing.T) {
	// Use a server that delays, then cancel the context
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	emb := newTestEmbedding(server.URL)
	require.NoError(t, emb.Start())

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	output, err := emb.Embed(ctx, agent.EmbeddingInput{Text: []string{"test"}})
	assert.Error(t, err)
	assert.Empty(t, output.Embedding)
}

func TestEmbedContextDeadline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	emb := newTestEmbedding(server.URL)
	require.NoError(t, emb.Start())

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Hour)) // already expired
	defer cancel()

	output, err := emb.Embed(ctx, agent.EmbeddingInput{Text: []string{"test"}})
	assert.Error(t, err)
	assert.Empty(t, output.Embedding)
}

// Result ordering validation
func TestEmbedResultOrdering(t *testing.T) {
	inputs := []string{"first", "second", "third"}
	expected := [][]float64{
		{1.0, 0.0, 0.0},
		{0.0, 1.0, 0.0},
		{0.0, 0.0, 1.0},
	}
	server := mockEmbeddingServer(t, expected, http.StatusOK)
	emb := newTestEmbedding(server.URL)
	require.NoError(t, emb.Start())

	ctx := context.Background()
	output, err := emb.Embed(ctx, agent.EmbeddingInput{Text: inputs})
	require.NoError(t, err)
	assert.Len(t, output.Embedding, len(inputs))
	for i := range inputs {
		assert.Equal(t, expected[i], output.Embedding[i],
			"embedding[%d] should correspond to input[%d] (%q)", i, i, inputs[i])
	}
}

// Concurrent safety after start
func TestEmbedConcurrentCalls(t *testing.T) {
	server := mockEmbeddingServer(t, [][]float64{{0.1, 0.2}}, http.StatusOK)
	emb := newTestEmbedding(server.URL)
	require.NoError(t, emb.Start())

	ctx := context.Background()
	done := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go func() {
			output, err := emb.Embed(ctx, agent.EmbeddingInput{Text: []string{"concurrent"}})
			assert.NoError(t, err)
			assert.Len(t, output.Embedding, 1)
			done <- true
		}()
	}

	for i := 0; i < 5; i++ {
		<-done
	}
}

// NewEmbedding factory default case
func TestNewEmbeddingDefaultType(t *testing.T) {
	config := EmbeddingConfig{
		Type:  ModelType(99), // unknown type
		Model: "text-embedding-3-small",
	}
	emb := NewEmbedding(config)
	require.NotNil(t, emb)
	_, ok := emb.(*openAiEmbedding)
	assert.True(t, ok, "unknown type should default to openAiEmbedding")
}
