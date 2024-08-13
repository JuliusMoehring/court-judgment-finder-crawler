package pdf

import "context"

type Reader interface {
	Read(ctx context.Context, path string) ([]byte, error)
}
