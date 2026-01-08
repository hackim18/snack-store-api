package config

import (
	"github.com/gin-gonic/gin"
)

func NewGin() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.SetTrustedProxies(nil)
	return engine
}
