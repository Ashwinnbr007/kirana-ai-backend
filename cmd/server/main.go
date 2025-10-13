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

	s3Store, err := storage.NewS3Storage(cfg.AWSConfig.S3Bucket, log, cfg.App.LocalDir)
	if err != nil {
		logger.L().Fatal("failed to init s3 storage", zap.Error(err))
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

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.App.Database.Host, cfg.App.Database.Port, cfg.App.Database.User, os.Getenv(models.DATABASE_PASSWORD_KEY), cfg.App.Database.DbName, "disable")

	db := storage.NewReposiory(cfg.App.Database.Type, connStr)
	if db == nil {
		zap.L().Panic("Could not initialize database")
	}

	// Initialise services
	audioService := service.NewAudioService(s3Store)
	aiService := service.NewAiService(openAiClient, restyClient, db)

	// Initialise handlers
	audioHandler := httpadapter.NewAudioHandler(audioService)
	aiHandler := httpadapter.NewAiHandler(aiService)

	v1 := router.Group("/v1")
	{
		v1.POST("/upload", audioHandler.UploadAudio)

		v1.POST("/transcribe/:fileName", aiHandler.TranscribeAudio)
		v1.POST("/translate_to_english", aiHandler.TranslateToEnglish)
		v1.POST("/english_to_inventory", aiHandler.InventoryDataToJsonTranslation)
		v1.POST("/english_to_sales", aiHandler.SalesDataToJsonTranslation)
	}

	router.Run(fmt.Sprintf(":%d", cfg.App.Port))
}
