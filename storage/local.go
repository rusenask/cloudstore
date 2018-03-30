package storage

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	Storage
	directory string
}

func NewLocalStorage(directory string) *LocalStorage {
	return &LocalStorage{
		directory: directory,
	}
}

func (ls *LocalStorage) Setup() error {
	return os.MkdirAll(ls.directory, os.ModePerm)
}

func (ls *LocalStorage) PublicURL(filename string) string {
	return ls.directory + "/" + filename
}

func (ls *LocalStorage) Store(ctx context.Context, filename string, data []byte, metadata map[string]string) error {
	return ioutil.WriteFile(filepath.Join(ls.directory, filename), data, 0644)
}

func (ls *LocalStorage) Get(ctx context.Context, filename string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(ls.directory, filename))
}

func (ls *LocalStorage) Delete(ctx context.Context, filename string) error {
	return os.Remove(filepath.Join(ls.directory, filename))
}
