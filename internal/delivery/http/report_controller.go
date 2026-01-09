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

type ReportController struct {
	Log      *logrus.Logger
	UseCase  *usecase.ReportUseCase
	Validate *validator.Validate
}

func NewReportController(
	useCase *usecase.ReportUseCase,
	logger *logrus.Logger,
	validate *validator.Validate,
) *ReportController {
	return &ReportController{
		Log:      logger,
		UseCase:  useCase,
		Validate: validate,
	}
}

func (c *ReportController) Transactions(ctx *gin.Context) {
	request := new(model.ReportTransactionsRequest)
	request.Start = strings.TrimSpace(ctx.Query("start"))
	request.End = strings.TrimSpace(ctx.Query("end"))

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Validation failed : %+v", err)
		message := utils.TranslateValidationError(c.Validate, err)
		utils.HandleHTTPError(ctx, utils.Error(message, http.StatusBadRequest, err))
		return
	}

	response, err := c.UseCase.Transactions(ctx.Request.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to get report : %+v", err)
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(messages.ReportFetched, response)
	ctx.JSON(http.StatusOK, res)
}
