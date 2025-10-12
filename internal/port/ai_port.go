package port

import "context"

type AiPort interface {
	TranslateToEnglish(ctx context.Context, transcription, transcriptionLangugae string) error
	DataToJsonTranslation(ctx context.Context, prompt string) error
}
