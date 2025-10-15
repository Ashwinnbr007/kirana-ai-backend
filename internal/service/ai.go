package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Ashwinnbr007/kirana-ai-backend/internal/constants"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/models"
	promptfactory "github.com/Ashwinnbr007/kirana-ai-backend/internal/prompt_factory"
	"github.com/go-resty/resty/v2"
	openai "github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type DatabasePort interface {
	WriteInventoryData(ctx context.Context, data *models.InventoryData) error
	WriteMultipleInventoryData(ctx context.Context, data *[]models.InventoryData) error
	WriteSalesData(ctx context.Context, data *[]models.SalesData) error
	ReadInventoryData(ctx context.Context) ([]models.InventoryData, error)
}

type AiService struct {
	openAiClient *openai.Client
	restyClient  *resty.Client
	db           DatabasePort
}

func NewAiService(openAiClient *openai.Client, restyClient *resty.Client, db DatabasePort) *AiService {
	return &AiService{
		openAiClient: openAiClient,
		restyClient:  restyClient,
		db:           db,
	}
}

func (s *AiService) TranslateToEnglish(ctx context.Context, transcription, transcriptionLangugae string) (*openai.ChatCompletionResponse, error) {

	contextualData, err := s.db.ReadInventoryData(ctx)
	if err != nil {
		zap.L().Warn("there was an error trying to fetch the context for translation", zap.Any("error: ", err))
	}

	var sb strings.Builder
	sb.WriteString("Product Inventory:\n")

	for _, p := range contextualData {
		line := fmt.Sprintf("  - Item: %s, Qty: %.2f, Price: %.2f, Unit: %sf\n", p.Item, p.Quantity, p.TotalCostOfProduct, p.Unit)
		sb.WriteString(line)
	}

	translationPrompt := strings.Replace(
		promptfactory.TRANSLATION_PROMPT,
		promptfactory.CONTEXT_PLACEHOLDER,
		sb.String(),
		1,
	)

	resp, err := s.openAiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT5ChatLatest,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleDeveloper,
					Content: promptfactory.TRANSLATION_DEVELOPER,
				},
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: translationPrompt,
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

func (s *AiService) DataToJsonTranslation(ctx context.Context, chatText string, typeOfRecord string) (any, error) {

	var prompt string
	var jsonData any

	switch typeOfRecord {
	case constants.INVENTORY_RECORD_IDENTIFIER:
		prompt = promptfactory.INVENTORY_DATA_TO_JSON_PROMPT
		jsonData = &[]models.InventoryData{}
	case constants.SALES_RECORD_IDENTIFIER:
		prompt = promptfactory.SALES_DATA_TO_JSON_PROMPT
		jsonData = &[]models.SalesData{}
	default:
		return nil, fmt.Errorf("given type of record is neither sales nor inventory, please check and retry")
	}

	resp, err := s.openAiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:           openai.GPT5Nano,
			ReasoningEffort: "minimal",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleDeveloper,
					Content: promptfactory.DATA_TO_JSON_DEVELOPER,
				},
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
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

	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &jsonData)
	if err != nil {
		zap.L().Error("error during unmarshaling of inventory data")
		return nil, fmt.Errorf("error during unmarshaling of inventory data")
	}

	switch typeOfRecord {
	case constants.INVENTORY_RECORD_IDENTIFIER:
		err = s.db.WriteMultipleInventoryData(ctx, jsonData.(*[]models.InventoryData))
	case constants.SALES_RECORD_IDENTIFIER:
		err = s.db.WriteSalesData(ctx, jsonData.(*[]models.SalesData))
	}

	if err != nil {
		zap.L().Error("error during updation of database")
		return nil, fmt.Errorf("error during updation of database")
	}

	return jsonData, nil
}

func (s *AiService) TranscribeAudio(ctx context.Context, audioFile string) (*models.TranscriptionResponse, error) {

	var transcriptionResponse *models.TranscriptionResponse

	resp, err := s.restyClient.R().
		SetContext(ctx).
		SetHeader("xi-api-key", os.Getenv("ELEVENLABS_API_KEY")).
		SetFile("file", audioFile).
		SetFormData(map[string]string{
			"model_id":      "scribe_v1",
			"language_code": "mal",
		}).
		SetResult(&transcriptionResponse).
		Post(constants.ELEVEN_LABS_BASE_URL + constants.ELEVEN_LABS_V1 + constants.ELEVEN_LABS_SPEECH_TO_TEXT_ENDPOINT)

	if err != nil {
		return nil, fmt.Errorf("unknonw error occurred while transcribing: %w", err)
	}

	if resp.IsError() {
		zap.L().Error("Request failed: ", zap.Any("status code: ", resp.StatusCode()))
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode())
	}

	return transcriptionResponse, nil
}
