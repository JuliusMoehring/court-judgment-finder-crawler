package filestorage

import "context"

type FileStorage interface {
	Exists(ctx context.Context, path string) (bool, error)
	Save(ctx context.Context, data []byte, path string) error
}
