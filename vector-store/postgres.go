package vectorstore

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/JuliusMoehring/court-judgment-finder-crawler/logger"
	"github.com/JuliusMoehring/court-judgment-finder-crawler/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

func uuidToString(uuid pgtype.UUID) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid.Bytes[0:4], uuid.Bytes[4:6], uuid.Bytes[6:8], uuid.Bytes[8:10], uuid.Bytes[10:16])
}

func getConfig() *pgxpool.Config {
	config, err := pgxpool.ParseConfig(os.Getenv("POSTGRES_CONNECTION_STRING"))
	if err != nil {
		panic(fmt.Sprintf("Failed to create a config, error: %s", err))
	}

	config.MaxConns = int32(12)
	config.MinConns = int32(4)
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30
	config.HealthCheckPeriod = time.Minute
	config.ConnConfig.ConnectTimeout = time.Second * 5

	return config
}

type PostgresVectorStore struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries

	logger logger.Logger
}

func NewPostgresVectorStore(ctx context.Context, logger logger.Logger) VectorStore {
	pool, err := pgxpool.NewWithConfig(ctx, getConfig())
	if err != nil {
		panic(err)
	}

	return &PostgresVectorStore{
		pool:    pool,
		queries: sqlc.New(pool),

		logger: logger,
	}
}

func (v *PostgresVectorStore) Close() {
	v.pool.Close()
}

func (v *PostgresVectorStore) CreateDocument(ctx context.Context, params CreateDocumentParams) error {
	tx, err := v.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queries := v.queries.WithTx(tx)

	documentID, err := queries.CreateDocument(ctx, params.FilePath)
	if err != nil {
		return err
	}

	for i, page := range params.Pages {
		_, err := queries.CreateDocumentPage(ctx, sqlc.CreateDocumentPageParams{
			Page:       int32(i + 1),
			Text:       page.Text,
			Embeddings: pgvector.NewVector(page.Embedding),
			DocumentID: documentID,
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (v *PostgresVectorStore) GetDocumentIDByFilePath(ctx context.Context, path string) (string, error) {
	uuid, err := v.queries.GetDocumentIDByFilePath(ctx, path)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return "", ErrDocumentNotFound
	}

	if err != nil {
		return "", err
	}

	return uuidToString(uuid), nil
}
