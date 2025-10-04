package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Ashwinnbr007/kirana-ai-backend/internal/models"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/pkg/logger"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/port"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

type AudioService struct {
	storage    port.StoragePort
	transcribe port.TranscriptionPort
}

func NewAudioService(storage port.StoragePort, transcribe port.TranscriptionPort) *AudioService {
	return &AudioService{storage: storage, transcribe: transcribe}
}

func (s *AudioService) SaveAudio(ctx context.Context, filename string, data []byte) (string, error) {

	finalName := fmt.Sprintf("%s_%s", time.Now().Format("20060102_150405"), filename)

	if err := s.storage.Save(ctx, finalName, data); err != nil {
		return "", err
	}

	return finalName, nil
}

func (s *AudioService) TranscribeAudio(ctx context.Context, jobName, bucket, key string) error {
	transcriptionJob, err := s.transcribe.GetTranscriptionJob(ctx, jobName, bucket, key)
	if err != nil {
		return fmt.Errorf("something went wrong during the transcription job: %w", err)
	}

	err = s.transcribe.GetTranscriptionResult(ctx, transcriptionJob)
	if err != nil {
		return fmt.Errorf("something went wrong fetching the result of the transcription job: %w", err)
	}

	return nil
}

func (a *AudioService) FetchTranscriptionJSON(ctx context.Context, bucket, key string) (*models.TranscriptionResponse, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket + "-transcription"), // hard coding `transcription` becasue having a config for each is a pain in the ass
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	var transcription models.TranscriptionResponse
	if err := json.NewDecoder(result.Body).Decode(&transcription); err != nil {
		logger.L().Error("failed to decode transcription JSON", zap.Error(err))
		return nil, err
	}

	return &transcription, nil
}
