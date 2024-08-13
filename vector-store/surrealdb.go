package vectorstore

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/conn/gorilla"
	"github.com/surrealdb/surrealdb.go/pkg/marshal"
)

type SurrealDBVectorStore struct {
	mu sync.Mutex

	db *surrealdb.DB
}

func NewSurrealDBVectorStore() VectorStore {
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
		db: db,
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
	defer v.db.Query("CANCEL TRANSACTION;", nil)

	var pageIDs []string

	for i, p := range params.Pages {
		pages, err := marshal.SmartUnmarshal[page](v.db.Create("page", map[string]interface{}{
			"page":      i + 1,
			"text":      p.Text,
			"embedding": p.Embedding,
		}))
		if err != nil {
			return err
		}

		pageIDs = append(pageIDs, pages[0].ID)
	}

	if _, err := v.db.Create("document", map[string]interface{}{
		"filePath": params.FilePath,
		"pages":    pageIDs,
	}); err != nil {
		return err
	}

	v.db.Query("COMMIT TRANSACTION;", nil)

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
