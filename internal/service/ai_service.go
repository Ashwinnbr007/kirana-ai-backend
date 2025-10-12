package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Ashwinnbr007/kirana-ai-backend/internal/models"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/port"
	promptfactory "github.com/Ashwinnbr007/kirana-ai-backend/internal/prompt_factory"
	openai "github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type AiService struct {
	aiPort       port.AiPort
	openAiClient *openai.Client
}

func NewAiService(aiPort port.AiPort, openAiClient *openai.Client) *AiService {
	return &AiService{aiPort: aiPort, openAiClient: openAiClient}
}

func (s *AiService) TranslateToEnglish(ctx context.Context, transcription, transcriptionLangugae string) (*openai.ChatCompletionResponse, error) {

	resp, err := s.openAiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT5ChatLatest,
			// ReasoningEffort: "minimal",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleDeveloper,
					Content: promptfactory.TRANSLATION_DEVELOPER,
				},
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: promptfactory.TRANSLATION_PROMPT,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: transcription,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("ChatCompletion error: %w", err)
	}

	if strings.Contains(resp.Choices[0].Message.Content, "Sorry, please input inventory or sales data only!") {
		return nil, fmt.Errorf("sorry, please input inventory or sales data only")
	}

	return &resp, nil
}

func (s *AiService) DataToJsonTranslation(ctx context.Context, chatText string) (*[]models.InventoryData, error) {

	resp, err := s.openAiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:           openai.GPT5Mini,
			ReasoningEffort: "minimal",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleDeveloper,
					Content: promptfactory.DATA_TO_JSON_DEVELOPER,
				},
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: promptfactory.DATA_TO_JSON_PROMPT,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: chatText,
				},
			},
		},
	)

	if err != nil {
		zap.L().Error("ChatCompletion error: %w")
		return nil, fmt.Errorf("ChatCompletion error: %w", err)
	}

	var inventoryData *[]models.InventoryData

	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &inventoryData)
	if err != nil {
		zap.L().Error("error during unmarshaling of inventory data")
		return nil, fmt.Errorf("error during unmarshaling of inventory data")
	}

	return inventoryData, nil
}
