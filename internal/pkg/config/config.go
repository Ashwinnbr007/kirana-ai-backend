package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AWSConfig struct {
	UseS3     bool   `mapstructure:"use_s3"`
	S3Bucket  string `mapstructure:"bucket"`
	AWSRegion string `mapstructure:"region"`
}

type Config struct {
	AppPort   int       `mapstructure:"app_port" validate:"required"`
	AWSConfig AWSConfig `mapstructure:"aws_config" validate:"omitempty"`
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
