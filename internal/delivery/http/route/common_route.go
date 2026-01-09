package route

import (
	"net/http"
	"snack-store-api/internal/messages"
	"snack-store-api/internal/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (c *RouteConfig) RegisterCommonRoutes(app *gin.Engine) {
	welcomeHandler := func(ctx *gin.Context) {
		payload := gin.H{"status": "ok"}
		res := utils.SuccessResponse(messages.WelcomeMessage, payload)
		ctx.JSON(http.StatusOK, res)
	}

	healthHandler := func(ctx *gin.Context) {
		payload := gin.H{"status": "ok"}
		res := utils.SuccessResponse(messages.HealthCheckSuccess, payload)
		ctx.JSON(http.StatusOK, res)
	}

	app.GET("/", welcomeHandler)
	app.GET("/api", welcomeHandler)
	app.GET("/health", healthHandler)
	app.StaticFile("/api/openapi.yaml", "api/openapi.yaml")
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/openapi.yaml")))
	app.NoRoute(func(ctx *gin.Context) {
		utils.HandleHTTPError(ctx, utils.Error(messages.NotFound, http.StatusNotFound, nil))
	})
}
