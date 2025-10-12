package httpadapter

import (
	"fmt"
	"strings"

	"github.com/Ashwinnbr007/kirana-ai-backend/internal/models"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/pkg/config"
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
		var apiError models.APIResponse

		if strings.Contains(err.Error(), "sorry, please input inventory or sales data only") {
			apiError = models.APIResponse{
				Status:  models.ErrInvalidInput,
				Message: "Error, please check your input",
				Data:    err.Error(),
			}
		} else {
			apiError = models.APIResponse{
				Status:  models.ErrInternal,
				Message: "An unknown error occurred",
				Data:    err.Error(),
			}
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

func (a *AiHandler) DataToJsonTranslation(c *gin.Context) {
	var inventoryApiBody models.InventoryApiBody
	if err := c.BindJSON(&inventoryApiBody); err != nil {
		apiError := models.APIResponse{
			Status:  models.ErrInvalidInput,
			Message: "Bad request, the request is invalid",
		}
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}

	ctx := c.Request.Context()
	var inventoryData *[]models.InventoryData
	inventoryData, err := a.aiService.DataToJsonTranslation(ctx, inventoryApiBody.InventoryInput)
	if err != nil {
		apiError := models.APIResponse{
			Status:  models.ErrInternal,
			Message: "An unknown error occurred",
		}
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}

	apiResponse := models.APIResponse{
		Status:  models.StatusOK,
		Message: "successfully converted to inventory data",
		Data:    inventoryData,
	}
	c.JSON(apiResponse.ToHTTPStatus(), apiResponse)
}

func (a *AiHandler) TranscribeAudio(c *gin.Context) {
	fileName := c.Param("fileName")
	ctx := c.Request.Context()

	if fileName == "" {
		apiError := models.APIResponse{
			Status:  models.ErrInvalidInput,
			Message: "file name is not provided in the path",
		}
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}

	projectRoot, _ := config.FindProjectRoot()
	filePath := fmt.Sprintf("%s/%s/%s", projectRoot, models.LOCAL_SAVE_PATH, fileName)

	transcriptionResponse, err := a.aiService.TranscribeAudio(ctx, filePath)
	if err != nil {
		apiError := models.APIResponse{
			Status:  models.ErrInternal,
			Message: "An unknown error occurred",
			Data:    err.Error(),
		}
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}

	apiResponse := models.APIResponse{
		Status:  models.StatusOK,
		Message: "successfully transcribed audio",
		Data:    transcriptionResponse,
	}
	c.JSON(apiResponse.ToHTTPStatus(), apiResponse)
}
