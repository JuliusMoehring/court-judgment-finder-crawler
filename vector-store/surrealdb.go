package vectorstore

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/JuliusMoehring/court-judgment-finder-crawler/logger"
	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/conn/gorilla"
	"github.com/surrealdb/surrealdb.go/pkg/marshal"
)

type SurrealDBVectorStore struct {
	mu sync.Mutex

	logger logger.Logger
	db     *surrealdb.DB
}

func NewSurrealDBVectorStore(logger logger.Logger) VectorStore {
	ws, err := gorilla.Create().Connect("ws://localhost:8000/rpc")
	if err != nil {
		panic(err)
	}

	db, err := surrealdb.New(os.Getenv("SURREAL_CONNECTION_STRING"), ws)
	if err != nil {
		panic(err)
	}

	if _, err = db.Signin(&surrealdb.Auth{
		Username: os.Getenv("SURREAL_USER"),
		Password: os.Getenv("SURREAL_PASSWORD"),
	}); err != nil {
		panic(err)
	}

	if _, err = db.Use(os.Getenv("SURREAL_NAMESPACE"), os.Getenv("SURREAL_DATABASE")); err != nil {
		panic(err)
	}

	return &SurrealDBVectorStore{
		logger: logger,
		db:     db,
	}
}

func (v *SurrealDBVectorStore) Close() {
	v.db.Close()
}

func (v *SurrealDBVectorStore) CreateDocument(ctx context.Context, params CreateDocumentParams) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	type page struct {
		ID        string    `json:"id,omitempty"`
		Page      int       `json:"page"`
		Text      string    `json:"text"`
		Embedding []float32 `json:"embedding"`
		CreatedAt time.Time `json:"createdAt,omitempty"`
		UpdatedAt time.Time `json:"updatedAt,omitempty"`
	}

	v.db.Query("BEGIN TRANSACTION;", nil)

	var pageIDs []string

	for i, p := range params.Pages {
		pages, err := marshal.SmartUnmarshal[page](v.db.Create("page", map[string]interface{}{
			"page":      i + 1,
			"text":      p.Text,
			"embedding": p.Embedding,
		}))
		if err != nil || len(pages) != 1 {
			v.logger.Errorf("vector-store", "failed to create page %d for path '%s'. Cancelling transaction.", i+1, params.FilePath)
			v.db.Query("CANCEL TRANSACTION;", nil)
			return err
		}

		pageIDs = append(pageIDs, pages[0].ID)
	}

	if _, err := v.db.Create("document", map[string]interface{}{
		"filePath": params.FilePath,
		"pages":    pageIDs,
	}); err != nil {
		v.logger.Errorf("vector-store", "failed to create document for path '%s'. Cancelling transaction.", params.FilePath)
		v.db.Query("CANCEL TRANSACTION;", nil)
		return err
	}

	if _, err := v.db.Query("COMMIT TRANSACTION;", nil); err != nil {
		v.logger.Errorf("vector-store", "failed to commit transaction. Cancelling transaction.")
		v.db.Query("CANCEL TRANSACTION;", nil)
		return err
	}

	return nil
}

var ErrDocumentNotFound = errors.New("document not found")

func (v *SurrealDBVectorStore) GetDocumentIDByPath(ctx context.Context, path string) (string, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	type result struct {
		ID string `json:"id"`
	}

	ids, err := marshal.SmartUnmarshal[result](v.db.Query("SELECT id FROM document WHERE filePath = $path;", map[string]string{"path": path}))
	if err != nil {
		return "", err
	}

	if len(ids) != 1 {
		return "", ErrDocumentNotFound
	}

	return ids[0].ID, nil
}
