package route

import (
	"snack-store-api/internal/delivery/http"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	Router                *gin.Engine
	CustomerController    *http.CustomerController
	ProductController     *http.ProductController
	TransactionController *http.TransactionController
	RedemptionController  *http.RedemptionController
	ReportController      *http.ReportController
	RateLimiter           gin.HandlerFunc
}

func (c *RouteConfig) Setup() {
	api := c.Router.Group("/api")
	api.Use(c.RateLimiter)

	c.RegisterCustomerRoutes(api)
	c.RegisterProductRoutes(api)
	c.RegisterTransactionRoutes(api)
	c.RegisterRedemptionRoutes(api)
	c.RegisterReportRoutes(api)
	c.RegisterCommonRoutes(c.Router)
}
