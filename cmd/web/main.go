package main

import (
	"fmt"
	"snack-store-api/internal/cache"
	"snack-store-api/internal/command"
	"snack-store-api/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	redisClient := config.NewRedis(viperConfig)
	cacheClient := cache.NewRedisCache(redisClient)
	executor := command.NewCommandExecutor(viperConfig, db)
	validate := config.NewValidator()
	router := config.NewGin()

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		Router:   router,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
		Cache:    cacheClient,
	})

	if !executor.Execute(log) {
		return
	}

	webPort := viperConfig.GetInt("PORT")
	err := router.Run(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
