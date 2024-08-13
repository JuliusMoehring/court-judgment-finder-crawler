package filestorage

import (
	"context"
	"os"
)

type LocalFileStorage struct {
}

func NewLocalFileStorage() FileStorage {
	return &LocalFileStorage{}
}

func (d *LocalFileStorage) Exists(ctx context.Context, path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	return false, nil
}

func (d *LocalFileStorage) Save(ctx context.Context, data []byte, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}
