package config

import (
	"snack-store-api/internal/cache"
	"snack-store-api/internal/delivery/http"
	"snack-store-api/internal/delivery/http/route"
	"snack-store-api/internal/repository"
	"snack-store-api/internal/usecase"

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
	Cache    cache.Cache
}

func Bootstrap(config *BootstrapConfig) {
	// Setup repositories
	customerRepository := repository.NewCustomerRepository(config.Log)
	productRepository := repository.NewProductRepository(config.Log)

	// Setup use cases
	customerUseCase := usecase.NewCustomerUseCase(
		config.DB,
		config.Log,
		customerRepository,
	)
	productUseCase := usecase.NewProductUseCase(
		config.DB,
		config.Log,
		productRepository,
		config.Cache,
	)

	// Setup controllers
	customerController := http.NewCustomerController(
		customerUseCase,
		config.Log,
		config.Validate,
	)
	productController := http.NewProductController(productUseCase, config.Log, config.Validate)

	// Setup routes
	routeConfig := route.RouteConfig{
		Router:             config.Router,
		CustomerController: customerController,
		ProductController:  productController,
	}
	routeConfig.Setup()
}
