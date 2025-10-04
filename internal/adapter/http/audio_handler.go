package httpadapter

import (
	"fmt"
	"io"

	"github.com/Ashwinnbr007/kirana-ai-backend/internal/models"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/pkg/config"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/pkg/logger"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AudioHandler struct {
	audioService *service.AudioService
}

func NewAudioHandler(audioService *service.AudioService) *AudioHandler {
	return &AudioHandler{audioService: audioService}
}

func (h *AudioHandler) UploadAndTranscribeAudio(c *gin.Context) {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		logger.L().Error("failed to load config: %w", zap.Error(err))
	}
	file, err := c.FormFile("file")
	if err != nil {
		apiError := models.APIResponse{
			Status:  models.ErrInvalidInput,
			Message: "file is required",
		}
		zap.L().Error("error, file not provided")
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}

	f, err := file.Open()
	if err != nil {
		apiError := models.APIResponse{
			Status:  models.ErrInternal,
			Message: "could not open file, please try again",
		}
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		apiError := models.APIResponse{
			Status:  models.ErrInternal,
			Message: "could not read file",
		}
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}
	ctx := c.Request.Context()
	fileName, err := h.audioService.SaveAudio(ctx, file.Filename, data)
	if err != nil {
		apiError := models.APIResponse{
			Status:  models.ErrInternal,
			Message: "could not save file",
		}
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}

	transcriptionJobName := fmt.Sprintf("%s_transcription_job", fileName)
	err = h.audioService.TranscribeAudio(ctx, transcriptionJobName, cfg.AWSConfig.S3Bucket, fileName)
	if err != nil {
		logger.L().Error("there was an error during transcription: %w", zap.Error(err))
	}
	transcriptionJobNameJSON := fmt.Sprintf("%s.json", transcriptionJobName)
	transcriptionResponse, err := h.audioService.FetchTranscriptionJSON(ctx, cfg.AWSConfig.S3Bucket, transcriptionJobNameJSON)
	if err != nil {
		logger.L().Error("an error occured while fetching the transcription: %w", zap.Error(err))
	}

	apiResponse := models.APIResponse{
		Status:  models.StatusCreated,
		Message: "successfully uploaded and transcribed",
		Data:    transcriptionResponse,
	}
	c.JSON(apiResponse.ToHTTPStatus(), apiResponse)
}
