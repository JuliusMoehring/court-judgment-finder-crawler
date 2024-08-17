package embedder

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/sashabaranov/go-openai"
)

var NoEmbeddingsReturnedError = errors.New("no embeddings returned")

type OpenAIEmbedder struct {
	mu sync.Mutex

	client *openai.Client
}

func NewOpenAIEmbedder() Embedder {
	return &OpenAIEmbedder{
		client: openai.NewClient(os.Getenv("OPENAI_API_KEY")),
	}
}

func (e *OpenAIEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	request := openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.AdaEmbeddingV2,
	}

	response, err := e.client.CreateEmbeddings(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, NoEmbeddingsReturnedError
	}

	return response.Data[0].Embedding, nil
}
