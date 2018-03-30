package storage

import (
	"context"
	"io"
)

type Storage interface {
	Setup() error
	PublicURL(filename string) string
	Store(ctx context.Context, filename string, data []byte, metadata map[string]string) error
	Get(ctx context.Context, filename string) (io.ReadCloser, error)
	Delete(ctx context.Context, filename string) error
}
