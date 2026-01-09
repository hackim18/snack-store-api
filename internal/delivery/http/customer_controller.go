package http

import (
	"net/http"

	"snack-store-api/internal/messages"
	"snack-store-api/internal/usecase"
	"snack-store-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type CustomerController struct {
	Log      *logrus.Logger
	UseCase  *usecase.CustomerUseCase
	Validate *validator.Validate
}

func NewCustomerController(
	useCase *usecase.CustomerUseCase,
	logger *logrus.Logger,
	validate *validator.Validate,
) *CustomerController {
	return &CustomerController{
		Log:      logger,
		UseCase:  useCase,
		Validate: validate,
	}
}

func (c *CustomerController) List(ctx *gin.Context) {
	response, err := c.UseCase.List(ctx.Request.Context())
	if err != nil {
		c.Log.Warnf("Failed to get customers : %+v", err)
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(messages.CustomersFetched, response)
	ctx.JSON(http.StatusOK, res)
}
