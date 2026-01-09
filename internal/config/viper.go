package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetDefault("APP_NAME", "snack-store-api")
	config.SetDefault("PORT", 8080)
	config.SetDefault("LOG_LEVEL", "info")
	config.SetDefault("DB_HOST", "localhost")
	config.SetDefault("DB_PORT", 5432)
	config.SetDefault("DB_NAME", "snack_store")
	config.SetDefault("DB_POOL_IDLE", 10)
	config.SetDefault("DB_POOL_MAX", 100)
	config.SetDefault("DB_POOL_LIFETIME", 300)
	config.SetDefault("REDIS_HOST", "localhost")
	config.SetDefault("REDIS_PORT", 6379)
	config.SetDefault("REDIS_PASSWORD", "")
	config.SetDefault("REDIS_DB", 0)

	config.SetConfigFile(".env")

	if err := config.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) {
			log.Println("No .env file found in root directory")
		} else {
			log.Printf("Error reading .env file: %v", err)
		}
	} else {
		log.Println("Successfully loaded configuration from .env")
	}

	config.AutomaticEnv()
	return config
}
