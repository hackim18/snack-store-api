package config

import (
	"snack-store-api/internal/delivery/http/route"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	Router   *gin.Engine
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	routeConfig := route.RouteConfig{
		Router: config.Router,
	}
	routeConfig.Setup()
}
