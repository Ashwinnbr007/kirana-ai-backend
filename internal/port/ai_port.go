package port

import (
	"context"

	"github.com/Ashwinnbr007/kirana-ai-backend/internal/models"
)

type AiPort interface {
	TranslateToEnglish(ctx context.Context, transcription, transcriptionLangugae string) error
	DataToJsonTranslation(ctx context.Context, prompt string) error
	TranscribeAudio(ctx context.Context, audioFile string) (models.TranscriptionResponse, error)
}
