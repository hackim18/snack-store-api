package http

import (
	"net/http"
	"strings"

	"snack-store-api/internal/messages"
	"snack-store-api/internal/model"
	"snack-store-api/internal/usecase"
	"snack-store-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type RedemptionController struct {
	Log      *logrus.Logger
	UseCase  *usecase.RedemptionUseCase
	Validate *validator.Validate
}

func NewRedemptionController(
	useCase *usecase.RedemptionUseCase,
	logger *logrus.Logger,
	validate *validator.Validate,
) *RedemptionController {
	return &RedemptionController{
		Log:      logger,
		UseCase:  useCase,
		Validate: validate,
	}
}

func (c *RedemptionController) Create(ctx *gin.Context) {
	request := new(model.CreateRedemptionRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		utils.HandleHTTPError(ctx, utils.Error(messages.FailedDataFromBody, http.StatusBadRequest, err))
		return
	}

	request.CustomerName = strings.TrimSpace(request.CustomerName)
	request.ProductID = strings.TrimSpace(request.ProductID)
	request.RedeemAt = strings.TrimSpace(request.RedeemAt)

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Validation failed : %+v", err)
		message := utils.TranslateValidationError(c.Validate, err)
		utils.HandleHTTPError(ctx, utils.Error(message, http.StatusBadRequest, err))
		return
	}

	response, err := c.UseCase.Create(ctx.Request.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to create redemption : %+v", err)
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(messages.RedemptionCreated, response)
	ctx.JSON(http.StatusCreated, res)
}
