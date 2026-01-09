package http

import (
	"net/http"

	"snack-store-api/internal/messages"
	"snack-store-api/internal/model"
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

const (
	defaultCustomerPage     = 1
	defaultCustomerPageSize = 10
)

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
	request := new(model.GetCustomerRequest)
	page, pageSize, err := utils.ParsePagination(
		ctx.Query("page"),
		ctx.Query("page_size"),
		defaultCustomerPage,
		defaultCustomerPageSize,
	)
	if err != nil {
		c.Log.Warnf("Failed to parse pagination : %+v", err)
		utils.HandleHTTPError(ctx, utils.Error(messages.FailedInputFormat, http.StatusBadRequest, err))
		return
	}

	request.Page = page
	request.PageSize = pageSize

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Validation failed : %+v", err)
		message := utils.TranslateValidationError(c.Validate, err)
		utils.HandleHTTPError(ctx, utils.Error(message, http.StatusBadRequest, err))
		return
	}

	response, paging, err := c.UseCase.List(ctx.Request.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to get customers : %+v", err)
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessWithPaginationResponse(messages.CustomersFetched, response, paging)
	ctx.JSON(http.StatusOK, res)
}
