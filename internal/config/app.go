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
	transactionRepository := repository.NewTransactionRepository(config.Log)
	redemptionRepository := repository.NewRedemptionRepository(config.Log)

	// Setup use cases
	customerUseCase := usecase.NewCustomerUseCase(config.DB, config.Log, customerRepository)
	productUseCase := usecase.NewProductUseCase(config.DB, config.Log, productRepository, config.Cache)
	transactionUseCase := usecase.NewTransactionUseCase(config.DB, config.Log, customerRepository, productRepository, transactionRepository)
	redemptionUseCase := usecase.NewRedemptionUseCase(config.DB, config.Log, customerRepository, productRepository, redemptionRepository)

	// Setup controllers
	customerController := http.NewCustomerController(customerUseCase, config.Log, config.Validate)
	productController := http.NewProductController(productUseCase, config.Log, config.Validate)
	transactionController := http.NewTransactionController(transactionUseCase, config.Log, config.Validate)
	redemptionController := http.NewRedemptionController(redemptionUseCase, config.Log, config.Validate)

	// Setup routes
	routeConfig := route.RouteConfig{
		Router:                config.Router,
		CustomerController:    customerController,
		ProductController:     productController,
		TransactionController: transactionController,
		RedemptionController:  redemptionController,
	}
	routeConfig.Setup()
}
