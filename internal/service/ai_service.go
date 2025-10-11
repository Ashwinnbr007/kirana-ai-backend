package service

import (
	"context"
	"fmt"

	"github.com/Ashwinnbr007/kirana-ai-backend/internal/pkg/config"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/pkg/logger"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/port"
	openai "github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type AiService struct {
	aiPort port.AiPort
}

func NewAiService(aiPort port.AiPort) *AiService {
	return &AiService{aiPort: aiPort}
}

func (s *AiService) TranslateToEnglish(ctx context.Context, transcription, transcriptionLangugae string) (string, error) {

	cfg, err := config.LoadConfig(".")
	if err != nil {
		logger.L().Error("failed to load config: %w", zap.Error(err))
	}

	client := openai.NewClient(cfg.AiConfiguaration.ApiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleDeveloper,
					Content: "Give clean English text output when you get Malayalam text",
				},
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `System Instructions:
You are a world-class language translator specializing in Malayalam-to-English translation, focused specifically on converting unstructured retail/inventory data.

Your primary goal is to **accurately translate and format inventory data** from the Malayalam text.

**Strict Output Rules:**
1.  **Product Names:** Do not interpret common product names (e.g., "വാട്ടർ മലൻ" must be "Water Melon").
2.  **Units & Pricing:** Simplify pricing units. Avoid redundant phrasing like "1 packet 15rs." Use the format [Quantity] [Price] [Unit Price] or [Quantity] [Total Price].
3.  **Format:** Each item must be on a new line and must include: [Product Name] [Quantity] [Price] [Unit].
4.  **Strict Filtering:** **ABSOLUTELY** ignore any irrelevant conversational filler, questions, or introductory/concluding remarks (e.g., "എത്രയരുന്നുള്ള"). Only capture saleable inventory data.

**Few-Shot Examples:**

Example 1:
Malayalam Input: ഹൈഡൻ സീക് ബിസ്കറ്റ് ഫിഫ്റ്റി കിലോ ഫൈവ് ഫിഫ്റ്റി റുപ്പീസ് ടോട്ടൽ പച്ചരി ത്രീ ഹൺഡ്രഡ് കിലോസ് ഫിഫ്റ്റി റുപ്പീ കിലോ പഴം ഫൈവ് ഹൺഡ്രഡ് കിലോസ് ട്വൻറ്റി റുപ്പീസ് പെർ കിലോ മൽബൊറോ സിഗരറ്റ് ഫൈവ് ഹൺഡ്രഡ് പാക്കറ്റ് പെർ പാക്കറ്റ് എത്രയരുന്നുള്ള തേർട്ടി റുപ്പീസ്
Output:
Hide and seek biscuit 50kg 550rs total
White rice 300kg 50rs per kilo
Banana 500kg 20rs per kilo
Marlboro Cigarette 500packet 30rs per packet

Example 2 (Addressing Unit & Name Errors):
Malayalam Input: കുക്കുംബർ ഇരുപത് കിലോ പത്ത് രൂപ പെർ കിലോ ഹൈഡൻസിക് ബിസ്കറ്റ് ടു ഹൺഡ്രഡ് പാക്കറ്റ് വൺ പാക്കറ്റ് ഫിഫ്റ്റീൻ റുപ്പീസ് വാട്ടർ മലൻ ഇരുന്നൂറ് കിലോ ടോട്ടൽ ഫൈവ് ഹൺഡ്രഡ് റുപ്പീസ്
Output:
Cucumber 20kg 10rs per kilo
Hide and seek biscuit 200packet 15rs per packet
Water Melon 200kg 500rs total
`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: transcription,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
