package httpadapter

import (
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/models"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AiHandler struct {
	aiService *service.AiService
}

func NewAiHandler(aiService *service.AiService) *AiHandler {
	return &AiHandler{
		aiService: aiService,
	}
}

func (a *AiHandler) TranslateToEnglish(c *gin.Context) {

	var translateApiBody models.TranslateApiBody
	if err := c.BindJSON(&translateApiBody); err != nil {
		apiError := models.APIResponse{
			Status:  models.ErrInvalidInput,
			Message: "Bad request, the request is invalid",
		}
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}

	// checking if the request contains any unsupported languages
	if !models.IsSupportedLanguage(translateApiBody.Language) {
		apiError := models.APIResponse{
			Status:  models.ErrInvalidInput,
			Message: "You have provided an unsupported language, please try any one of the following languages",
			Data: map[string][]string{
				"supported_languages": models.GetSupportedLanguages(),
			},
		}
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}

	ctx := c.Request.Context()
	translatedText, err := a.aiService.TranslateToEnglish(ctx, translateApiBody.TextToTranslate, translateApiBody.Language)
	if err != nil {
		zap.L().Error("an error occured tyring to translate to english: ", zap.Any("error: ", err))
		apiError := models.APIResponse{
			Status:  models.ErrInternal,
			Message: "An unknown error occurred",
			Data:    err,
		}
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}

	apiResponse := models.APIResponse{
		Status:  models.StatusOK,
		Message: "translation successfull",
		Data:    translatedText,
	}
	c.JSON(apiResponse.ToHTTPStatus(), apiResponse)
}
