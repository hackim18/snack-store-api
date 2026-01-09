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
}

func (c *RouteConfig) Setup() {
	api := c.Router.Group("/api")

	c.RegisterCustomerRoutes(api)
	c.RegisterProductRoutes(api)
	c.RegisterTransactionRoutes(api)
	c.RegisterCommonRoutes(c.Router)
}
