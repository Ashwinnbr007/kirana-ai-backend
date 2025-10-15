package service

import (
	"context"
	"fmt"
	"time"
)

type StoragePort interface {
	Save(ctx context.Context, filename string, data []byte) error
}

type AudioService struct {
	storage StoragePort
}

func NewAudioService(storage StoragePort) *AudioService {
	return &AudioService{storage: storage}
}

func (s *AudioService) SaveAudio(ctx context.Context, filename string, data []byte) (string, error) {

	finalName := fmt.Sprintf("%s_%s", time.Now().Format("20060102_150405"), filename)

	if err := s.storage.Save(ctx, finalName, data); err != nil {
		return "", err
	}

	return finalName, nil
}
