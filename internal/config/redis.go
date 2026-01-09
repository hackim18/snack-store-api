package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func NewRedis(viper *viper.Viper) *redis.Client {
	host := viper.GetString("REDIS_HOST")
	port := viper.GetInt("REDIS_PORT")
	password := viper.GetString("REDIS_PASSWORD")
	db := viper.GetInt("REDIS_DB")

	address := fmt.Sprintf("%s:%d", host, port)
	return redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
}
