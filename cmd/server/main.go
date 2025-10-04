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

	cfg, err := config.LoadConfig("../../internal/pkg/config")
	if err != nil {
		logger.L().Fatal("failed to load config", zap.Error(err))
	}
	var store port.StoragePort
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
	audioService := service.NewAudioService(store)
	audioHandler := httpadapter.NewAudioHandler(audioService)

	v1 := router.Group("/v1")
	{
		v1.POST("/audio", audioHandler.UploadAudio)
	}

	router.Run(fmt.Sprintf(":%d", cfg.App.Port))
}
