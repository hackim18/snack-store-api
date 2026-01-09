package route

import "github.com/gin-gonic/gin"

func (c *RouteConfig) RegisterReportRoutes(rg *gin.RouterGroup) {
	reports := rg.Group("/reports")

	reports.GET("/transactions", c.ReportController.Transactions)
}
