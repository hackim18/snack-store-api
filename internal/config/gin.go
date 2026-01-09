package config

import (
	"time"

	"snack-store-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewGin(logger *logrus.Logger) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestLogger(logger, 2*time.Second))
	engine.SetTrustedProxies(nil)
	return engine
}
