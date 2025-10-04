package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Ashwinnbr007/kinara-ai-backend/internal/port"
)

type AudioService struct {
	storage port.StoragePort
}

func NewAudioService(storage port.StoragePort) *AudioService {
	return &AudioService{storage: storage}
}

func (s *AudioService) SaveAudio(ctx context.Context, filename string, data []byte) (string, error) {
	
	finalName := fmt.Sprintf("%s_%s", time.Now().Format("20060102_150405"), filename)

	if err := s.storage.Save(ctx, finalName, data); err != nil {
		return "", err
	}

	return finalName, nil
}
