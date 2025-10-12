package storage

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/zap"
)

type S3Storage struct {
	client     *s3.Client
	bucketName string
	logger     *zap.Logger
	localDir   string
}

func NewS3Storage(bucketName string, logger *zap.Logger, localDir string) (*S3Storage, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	return &S3Storage{
		client:     client,
		bucketName: bucketName,
		logger:     logger,
		localDir:   localDir,
	}, nil
}

func (s *S3Storage) Save(ctx context.Context, filename string, data []byte) error {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(data),
		ACL:         types.ObjectCannedACLPrivate,
		ContentType: aws.String("audio/m4a"),
	}

	if err := os.MkdirAll(s.localDir, 0755); err != nil {
		return err
	}

	dst := filepath.Join(s.localDir, filename)
	os.WriteFile(dst, data, 0644)

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		s.logger.Error("failed to upload to S3", zap.String("file", filename), zap.Error(err))
		return fmt.Errorf("s3 upload error: %w", err)
	}
	s.logger.Info("uploaded file to S3", zap.String("file", filename))

	return nil
}
