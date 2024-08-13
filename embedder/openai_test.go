package embedder

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func Test_OpenAIEmbedder(t *testing.T) {
	t.Run("", func(t *testing.T) {
		err := godotenv.Load("../.env")
		if err != nil {
			t.Logf("error loading .env file: %v", err)
		}

		embedder := NewOpenAIEmbedder()

		embedding, err := embedder.Embed(context.Background(), "Hello, my dog is cute")
		assert.NoError(t, err, "Should not return an error")
		assert.Len(t, embedding, 1536, "Should return the correct number of embeddings")
	})
}
