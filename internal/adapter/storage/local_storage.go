package storage

import (
	"context"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	dir string
}

func NewLocalStorage(dir string) *LocalStorage {
	return &LocalStorage{dir: dir}
}

func (s *LocalStorage) Save(ctx context.Context, filename string, data []byte) error {
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return err
	}

	dst := filepath.Join(s.dir, filename)
	return os.WriteFile(dst, data, 0644)
}
