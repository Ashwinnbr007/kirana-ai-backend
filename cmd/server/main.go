package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"

	httpadapter "github.com/Ashwinnbr007/kirana-ai-backend/internal/adapter/http"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/adapter/storage"
	"github.com/Ashwinnbr007/kirana-ai-backend/internal/models"
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
	configPath := fmt.Sprintf("%s/%s", projectRoot, models.CONFIG_PATH)
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

	models.InitSupportedLanguages(cfg.App.SupportedLanguages)
	router := gin.Default()
	if err != nil {
		logger.L().Fatal("failed to init transcription storage", zap.Error(err))
	}
	// Initialise clients
	openaiApiKey := os.Getenv("OPENAI_API_KEY")
	if openaiApiKey == "" {
		logger.L().Fatal("please check your env variable: OPENAI_API_KEY, looks like it is null", zap.Error(err))
	}
	openAiClient := openai.NewClient(openaiApiKey)
	restyClient := resty.New()

	// Initialise services
	audioService := service.NewAudioService(store)
	aiService := service.NewAiService(aiPort, openAiClient, restyClient)

	// Initialise handlers
	audioHandler := httpadapter.NewAudioHandler(audioService)
	aiHandler := httpadapter.NewAiHandler(aiService)

	v1 := router.Group("/v1")
	{
		v1.POST("/upload", audioHandler.UploadAudio)

		v1.POST("/transcribe/:fileName", aiHandler.TranscribeAudio)
		v1.POST("/translate_to_english", aiHandler.TranslateToEnglish)
		v1.POST("/english_to_inventory", aiHandler.DataToJsonTranslation)
	}

	router.Run(fmt.Sprintf(":%d", cfg.App.Port))
}
