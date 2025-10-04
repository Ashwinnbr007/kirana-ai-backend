package main

import (
	"github.com/gin-gonic/gin"

	httpadapter "github.com/Ashwinnbr007/kinara-ai-backend/internal/adapter/http"
	"github.com/Ashwinnbr007/kinara-ai-backend/internal/adapter/storage"
	"github.com/Ashwinnbr007/kinara-ai-backend/internal/pkg/logger"

	"github.com/Ashwinnbr007/kinara-ai-backend/internal/service"
)

func main() {

	if err := logger.Init(); err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()
	log := logger.L()
	log.Info("Logger initialized successfully")

	router := gin.Default()

	storageAdapter := storage.NewLocalStorage("uploads")
	audioService := service.NewAudioService(storageAdapter)
	audioHandler := httpadapter.NewAudioHandler(audioService)

	v1 := router.Group("/v1")
	{
		v1.POST("/audio", audioHandler.UploadAudio)
	}

	router.Run(":8080")
}
