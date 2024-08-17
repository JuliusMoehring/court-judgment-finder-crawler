package vectorstore

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

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
		Page      int       `json:"page"`
		Text      string    `json:"text"`
		Embedding []float32 `json:"embedding"`
	}

	var pages []page

	for i, p := range params.Pages {
		pages = append(pages, page{
			Page:      i + 1,
			Text:      p.Text,
			Embedding: p.Embedding,
		})
	}

	response, err := v.db.Query(`
		BEGIN TRANSACTION;

		LET $doc = (CREATE ONLY document SET filePath = $filePath);

		INSERT INTO page (SELECT *, [$doc.id, page] AS id FROM $pages);

		COMMIT TRANSACTION;`,
		map[string]interface{}{
			"filePath": params.FilePath,
			"pages":    pages,
		})
	if err != nil {
		v.logger.Errorf("vector-store", "failed to create document for path '%s'.", params.FilePath)
		return err
	}

	var queryResult []marshal.RawQuery[any]

	if err := marshal.UnmarshalRaw(response, &queryResult); err != nil {
		v.logger.Errorf("vector-store", "failed to unmarshal response for path '%s': %s", params.FilePath, err)
		return err
	}

	for _, result := range queryResult {
		if result.Status != marshal.StatusOK {
			return errors.New(fmt.Sprintf("failed to create document for path '%s': %s", params.FilePath, result.Detail))
		}
	}

	return nil
}

func (v *SurrealDBVectorStore) GetDocumentIDByFilePath(ctx context.Context, path string) (string, error) {
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
