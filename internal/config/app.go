package config

import (
	"snack-store-api/internal/cache"
	"snack-store-api/internal/delivery/http"
	"snack-store-api/internal/delivery/http/middleware"
	"snack-store-api/internal/delivery/http/route"
	"snack-store-api/internal/repository"
	"snack-store-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	Router   *gin.Engine
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate
	Viper    *viper.Viper
	Cache    cache.Cache
	Redis    *redis.Client
}

func Bootstrap(config *BootstrapConfig) {
	// Setup repositories
	customerRepository := repository.NewCustomerRepository(config.Log)
	productRepository := repository.NewProductRepository(config.Log)
	transactionRepository := repository.NewTransactionRepository(config.Log)
	redemptionRepository := repository.NewRedemptionRepository(config.Log)
	reportRepository := repository.NewReportRepository(config.Log)

	// Setup use cases
	customerUseCase := usecase.NewCustomerUseCase(config.DB, config.Log, customerRepository)
	productUseCase := usecase.NewProductUseCase(config.DB, config.Log, productRepository, config.Cache)
	transactionUseCase := usecase.NewTransactionUseCase(config.DB, config.Log, customerRepository, productRepository, transactionRepository, config.Cache)
	redemptionUseCase := usecase.NewRedemptionUseCase(config.DB, config.Log, customerRepository, productRepository, redemptionRepository, config.Cache)
	reportUseCase := usecase.NewReportUseCase(config.DB, config.Log, reportRepository, config.Cache)

	// Setup controllers
	customerController := http.NewCustomerController(customerUseCase, config.Log, config.Validate)
	productController := http.NewProductController(productUseCase, config.Log, config.Validate)
	transactionController := http.NewTransactionController(transactionUseCase, config.Log, config.Validate)
	redemptionController := http.NewRedemptionController(redemptionUseCase, config.Log, config.Validate)
	reportController := http.NewReportController(reportUseCase, config.Log, config.Validate)

	// Setup middleware
	rateLimiterMiddleware := middleware.NewRateLimiter(config.Viper, config.Redis)

	// Setup routes
	routeConfig := route.RouteConfig{
		Router:                config.Router,
		CustomerController:    customerController,
		ProductController:     productController,
		TransactionController: transactionController,
		RedemptionController:  redemptionController,
		ReportController:      reportController,
		RateLimiter:           rateLimiterMiddleware,
	}
	routeConfig.Setup()
}
