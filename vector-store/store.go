package vectorstore

import "context"

type CreateDocumentParamsPage struct {
	Text      string
	Embedding []float32
}

type CreateDocumentParams struct {
	FilePath string
	Pages    []CreateDocumentParamsPage
}

type VectorStore interface {
	Close()

	CreateDocument(ctx context.Context, params CreateDocumentParams) error
	GetDocumentIDByPath(ctx context.Context, path string) (string, error)
}
