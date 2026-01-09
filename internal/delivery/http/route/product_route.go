package route

import "github.com/gin-gonic/gin"

func (c *RouteConfig) RegisterProductRoutes(rg *gin.RouterGroup) {
	products := rg.Group("/products")

	products.GET("", c.ProductController.ListByDate)
}
