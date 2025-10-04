package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/Ashwinnbr007/kirana-ai-backend/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/transcribe"
	"github.com/aws/aws-sdk-go-v2/service/transcribe/types"
	"go.uber.org/zap"
)

type Transcribe struct {
	client *transcribe.Client
	logger *zap.Logger
}

func NewTranscription(logger *zap.Logger) (*Transcribe, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := transcribe.NewFromConfig(cfg)
	return &Transcribe{client: client, logger: logger}, nil
}

func (t *Transcribe) GetTranscriptionJob(ctx context.Context, jobName, bucket, key string) (string, error) {

	language := models.MapLanguageToAWS(models.LanguageMalayalam)
	input := &transcribe.StartTranscriptionJobInput{
		TranscriptionJobName: aws.String(jobName),
		LanguageCode:         language,
		Media: &types.Media{
			MediaFileUri: aws.String(fmt.Sprintf("s3://%s/%s", bucket, key)),
		},
		OutputBucketName: aws.String("kirana-ai-audio-transcription"),
	}

	_, err := t.client.StartTranscriptionJob(ctx, input)
	if err != nil {
		t.logger.Error("failed to start transcription job", zap.String("job", jobName), zap.Error(err))
		return "", err
	}

	t.logger.Info("Transcription job started", zap.String("job_name", jobName))
	return jobName, nil
}

func (t *Transcribe) GetTranscriptionResult(ctx context.Context, jobName string) error {
	for {
		resp, err := t.client.GetTranscriptionJob(ctx, &transcribe.GetTranscriptionJobInput{
			TranscriptionJobName: aws.String(jobName),
		})
		if err != nil {
			return fmt.Errorf("failed to get transcription job: %w", err)
		}

		status := resp.TranscriptionJob.TranscriptionJobStatus
		t.logger.Info("Polling transcription job", zap.String("job_name", jobName), zap.String("status", string(status)))

		if status == types.TranscriptionJobStatusCompleted {
			t.logger.Info("Transcription completed", zap.String("uri", *resp.TranscriptionJob.Transcript.TranscriptFileUri))
			return nil
		}

		if status == types.TranscriptionJobStatusFailed {
			return fmt.Errorf("transcription job failed")
		}

		// Adding this so that there is no constant polling to the aws transcription server
		// Also allows the service to wait for transcription file to reach the s3 bucket as a JSON
		time.Sleep(5 * time.Second)
	}
}
