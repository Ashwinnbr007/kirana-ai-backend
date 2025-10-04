package port

import "context"

type TranscriptionPort interface {
	GetTranscriptionJob(ctx context.Context, jobName, bucket, key string) (string, error)
	GetTranscriptionResult(ctx context.Context, jobName string) error
}
