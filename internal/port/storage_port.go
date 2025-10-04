package port

import "context"

type StoragePort interface {
	Save(ctx context.Context, filename string, data []byte) error
}
