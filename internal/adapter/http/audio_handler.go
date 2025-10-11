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

func (h *AudioHandler) UploadAudio(c *gin.Context) {
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
	// Creating the final response
	responseData := map[string]string{
		"filename": fileName,
	}

	apiResponse := models.APIResponse{
		Status:  models.StatusCreated,
		Message: "file Uploaded succesfully",
		Data:    responseData,
	}
	c.JSON(apiResponse.ToHTTPStatus(), apiResponse)
}

func (h *AudioHandler) CreateTranscriptionJob(c *gin.Context) {

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

	cfg, err := config.LoadConfig(".")
	if err != nil {
		logger.L().Error("failed to load config: %w", zap.Error(err))
	}

	transcriptionJobName := fmt.Sprintf("%s_transcription_job", fileName)
	err = h.audioService.TranscribeAudio(ctx, transcriptionJobName, cfg.AWSConfig.S3Bucket, fileName)
	if err != nil {
		logger.L().Error("there was an error during creation of the transcription job: %w", zap.Error(err))
	}
	transcriptionJobNameJSON := fmt.Sprintf("%s.json", transcriptionJobName)

	// Creating the final response
	responseData := map[string]string{
		"transcrption_job_name": transcriptionJobNameJSON,
	}
	apiResponse := models.APIResponse{
		Status:  models.StatusAccepted,
		Message: "transcription job created",
		Data:    responseData,
	}
	c.JSON(apiResponse.ToHTTPStatus(), apiResponse)
}

func (h *AudioHandler) FetchTranscription(c *gin.Context) {

	fileName := c.Param("fileName")

	if fileName == "" {
		apiError := models.APIResponse{
			Status:  models.ErrInvalidInput,
			Message: "file name is not provided in the path",
		}
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}

	ctx := c.Request.Context()
	cfg, err := config.LoadConfig(".")
	if err != nil {
		logger.L().Error("failed to load config: %w", zap.Error(err))
	}

	transcriptionResponse, err := h.audioService.FetchTranscriptionJSON(ctx, cfg.AWSConfig.S3Bucket, fileName)
	if err != nil {
		logger.L().Error("there is some error tryng to fetch the transcrition: %w", zap.Error(err))
		apiError := models.APIResponse{
			Status:  models.ErrInternal,
			Message: "file name is not provided in the path",
		}
		c.JSON(apiError.ToHTTPStatus(), gin.H{"error": apiError})
		return
	}

	// Creating the final response
	apiResponse := models.APIResponse{
		Status:  models.StatusOK,
		Message: "successfully uploaded and transcribed",
		Data:    &transcriptionResponse.Results.Transcripts[0],
	}
	c.JSON(apiResponse.ToHTTPStatus(), apiResponse)
}
