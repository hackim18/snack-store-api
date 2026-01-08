package route

import "github.com/gin-gonic/gin"

type RouteConfig struct {
	Router *gin.Engine
}

func (c *RouteConfig) Setup() {
	c.RegisterCommonRoutes(c.Router)
}
