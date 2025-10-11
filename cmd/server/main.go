package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	httpadapter "github.com/Ashwinnbr007/kirana-ai-backend/internal/adapter/http"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/adapter/storage"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/pkg/config"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/pkg/logger"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/port"

	"github.com/Ashwinnbr007/kirana-ai-backend/internal/service"
)

func main() {
	if err := logger.Init(); err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()
	log := logger.L()
	log.Info("Logger initialized successfully")

	projectRoot, err := config.FindProjectRoot()
	if err != nil {
		logger.L().Error("could not find the project root", zap.Error(err))
		projectRoot = "."
	}
	configPath := fmt.Sprintf("%s/internal/pkg/config", projectRoot)
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.L().Fatal("failed to load config", zap.Error(err))
	}

	var store port.StoragePort
	var aiPort port.AiPort

	if cfg.AWSConfig.UseS3 {
		s3Store, err := storage.NewS3Storage(cfg.AWSConfig.S3Bucket, log)
		if err != nil {
			logger.L().Fatal("failed to init s3 storage", zap.Error(err))
		}
		store = s3Store
	} else {
		store = storage.NewLocalStorage("uploads")
	}

	router := gin.Default()
	transcriptionStore, err := storage.NewTranscription(log)
	if err != nil {
		logger.L().Fatal("failed to init transcription storage", zap.Error(err))
	}

	// Initialise services
	audioService := service.NewAudioService(store, transcriptionStore)
	aiService := service.NewAiService(aiPort)

	// Initialise handlers
	audioHandler := httpadapter.NewAudioHandler(audioService)
	aiHandler := httpadapter.NewAiHandler(aiService)

	v1 := router.Group("/v1")
	{
		v1.POST("/upload", audioHandler.UploadAudio)

		v1.POST("/transcribe/:fileName", audioHandler.CreateTranscriptionJob)
		v1.GET("/transcribe/:fileName", audioHandler.FetchTranscription)

		v1.POST("/translate", aiHandler.TranslateToEnglish)
	}

	router.Run(fmt.Sprintf(":%d", cfg.App.Port))
}
