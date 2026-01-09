package route

import (
	"snack-store-api/internal/delivery/http"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	Router            *gin.Engine
	ProductController *http.ProductController
}

func (c *RouteConfig) Setup() {
	api := c.Router.Group("/api")

	c.RegisterProductRoutes(api)
	c.RegisterCommonRoutes(c.Router)
}
