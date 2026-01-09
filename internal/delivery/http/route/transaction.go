package route

import "github.com/gin-gonic/gin"

func (c *RouteConfig) RegisterTransactionRoutes(rg *gin.RouterGroup) {
	transactions := rg.Group("/transactions")

	transactions.POST("", c.TransactionController.Create)
	transactions.GET("", c.TransactionController.List)
}
