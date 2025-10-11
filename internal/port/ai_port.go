package port

import "context"

type AiPort interface {
	TranslateToEnglish(ctx context.Context, prompt string) error
}
