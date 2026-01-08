package route

import (
	"net/http"
	"snack-store-api/internal/messages"

	"github.com/gin-gonic/gin"
)

func (c *RouteConfig) RegisterCommonRoutes(app *gin.Engine) {
	welcomeHandler := func(ctx *gin.Context) {
		res := gin.H{"message": messages.WelcomeMessage}
		ctx.JSON(http.StatusOK, res)
	}

	healthHandler := func(ctx *gin.Context) {
		res := gin.H{"status": "ok"}
		ctx.JSON(http.StatusOK, res)
	}

	app.GET("/", welcomeHandler)
	app.GET("/api", welcomeHandler)
	app.GET("/health", healthHandler)
	app.NoRoute(func(ctx *gin.Context) {
		res := gin.H{"message": messages.NotFound}
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
	})
}
