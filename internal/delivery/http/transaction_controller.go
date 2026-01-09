package http

import (
	"net/http"
	"strings"

	"snack-store-api/internal/constants"
	"snack-store-api/internal/messages"
	"snack-store-api/internal/model"
	"snack-store-api/internal/usecase"
	"snack-store-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type TransactionController struct {
	Log      *logrus.Logger
	UseCase  *usecase.TransactionUseCase
	Validate *validator.Validate
}

func NewTransactionController(
	useCase *usecase.TransactionUseCase,
	logger *logrus.Logger,
	validate *validator.Validate,
) *TransactionController {
	return &TransactionController{
		Log:      logger,
		UseCase:  useCase,
		Validate: validate,
	}
}

func (c *TransactionController) Create(ctx *gin.Context) {
	request := new(model.CreateTransactionRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		utils.HandleHTTPError(ctx, utils.Error(messages.FailedDataFromBody, http.StatusBadRequest, err))
		return
	}

	request.CustomerName = strings.TrimSpace(request.CustomerName)
	request.ProductID = strings.TrimSpace(request.ProductID)
	request.TransactionAt = strings.TrimSpace(request.TransactionAt)

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Validation failed : %+v", err)
		message := utils.TranslateValidationError(c.Validate, err)
		utils.HandleHTTPError(ctx, utils.Error(message, http.StatusBadRequest, err))
		return
	}

	response, err := c.UseCase.Create(ctx.Request.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to create transaction : %+v", err)
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(messages.TransactionCreated, response)
	ctx.JSON(http.StatusCreated, res)
}

func (c *TransactionController) List(ctx *gin.Context) {
	request := new(model.GetTransactionRequest)
	request.Start = strings.TrimSpace(ctx.Query("start"))
	request.End = strings.TrimSpace(ctx.Query("end"))
	page, pageSize, err := utils.ParsePagination(
		ctx.Query("page"),
		ctx.Query("page_size"),
		constants.DefaultPage,
		constants.DefaultPageSize,
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
		c.Log.Warnf("Failed to get transactions : %+v", err)
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessWithPaginationResponse(messages.TransactionsFetched, response, paging)
	ctx.JSON(http.StatusOK, res)
}
