package httpadapter

import (
	"context"
	"io"

	"github.com/Ashwinnbr007/kirana-ai-backend/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AudioService interface {
	SaveAudio(ctx context.Context, filename string, data []byte) (string, error)
}

type AudioHandler struct {
	audioService AudioService
}

func NewAudioHandler(audioService AudioService) *AudioHandler {
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
