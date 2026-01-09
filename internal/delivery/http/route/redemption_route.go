package route

import "github.com/gin-gonic/gin"

func (c *RouteConfig) RegisterRedemptionRoutes(rg *gin.RouterGroup) {
	redemptions := rg.Group("/redemptions")

	redemptions.POST("", c.RedemptionController.Create)
}
