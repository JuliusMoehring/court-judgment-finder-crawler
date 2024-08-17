package vectorstore

import (
	"context"
	"errors"
)

type CreateDocumentParamsPage struct {
	Text      string
	Embedding []float32
}

type CreateDocumentParams struct {
	FilePath string
	Pages    []CreateDocumentParamsPage
}

var ErrDocumentNotFound = errors.New("document not found")

type VectorStore interface {
	Close()

	CreateDocument(ctx context.Context, params CreateDocumentParams) error
	GetDocumentIDByFilePath(ctx context.Context, path string) (string, error)
}
