package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type AWSConfig struct {
	UseS3     bool   `mapstructure:"use_s3"`
	S3Bucket  string `mapstructure:"bucket"`
	AWSRegion string `mapstructure:"region"`
}

type App struct {
	Port int `mapstructure:"port" validate:"required"`
}

type Config struct {
	App       App       `mapstructure:"app" validate:"required"`
	AWSConfig AWSConfig `mapstructure:"aws_config" validate:"omitempty"`
}

func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("project root not found (no go.mod or .git)")
		}
		dir = parent
	}
}

func LoadConfig(path string) (Config, error) {

	var cfg Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return cfg, nil
}
