package filestorage

import (
	"bytes"
	"context"
	"os"

	"github.com/JuliusMoehring/court-judgment-finder-crawler/logger"
	"github.com/aws/smithy-go/ptr"
	supabase "github.com/supabase-community/storage-go"
)

type SupabaseFileStorage struct {
	client *supabase.Client
	bucket string

	logger logger.Logger
}

func NewSupabaseFileStorage(logger logger.Logger, bucket string) FileStorage {
	client := supabase.NewClient(os.Getenv("SUPABASE_STORAGE_URL"), os.Getenv("SUPABASE_PROJECT_SECRET_API_KEY"), nil)

	return &SupabaseFileStorage{
		client: client,
		bucket: bucket,
		logger: logger,
	}
}

func (s *SupabaseFileStorage) Exists(ctx context.Context, path string) (bool, error) {
	_, err := s.client.ListFiles(s.bucket, path, supabase.FileSearchOptions{})
	if err != nil {
		return false, err
	}

	return false, nil
}

func (s *SupabaseFileStorage) Save(ctx context.Context, data []byte, path string) error {
	_, err := s.client.UploadOrUpdateFile(s.bucket, path, bytes.NewReader(data), false, supabase.FileOptions{
		ContentType: ptr.String("application/pdf"),
	})
	if err != nil {
		return err
	}

	return nil
}
