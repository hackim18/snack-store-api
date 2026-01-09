package route

import "github.com/gin-gonic/gin"

func (c *RouteConfig) RegisterCustomerRoutes(rg *gin.RouterGroup) {
	customers := rg.Group("/customers")

	customers.GET("", c.CustomerController.List)
}
