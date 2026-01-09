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

type ProductController struct {
	Log      *logrus.Logger
	UseCase  *usecase.ProductUseCase
	Validate *validator.Validate
}

func NewProductController(
	useCase *usecase.ProductUseCase,
	logger *logrus.Logger,
	validate *validator.Validate,
) *ProductController {
	return &ProductController{
		Log:      logger,
		UseCase:  useCase,
		Validate: validate,
	}
}

func (c *ProductController) ListByDate(ctx *gin.Context) {
	request := new(model.GetProductRequest)
	request.Date = strings.TrimSpace(ctx.Query("date"))

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Validation failed : %+v", err)
		message := utils.TranslateValidationError(c.Validate, err)
		utils.HandleHTTPError(ctx, utils.Error(message, http.StatusBadRequest, err))
		return
	}

	response, err := c.UseCase.ListByDate(ctx.Request.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to get products : %+v", err)
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(messages.ProductsFetched, response)
	ctx.JSON(http.StatusOK, res)
}

func (c *ProductController) Create(ctx *gin.Context) {
	request := new(model.CreateProductRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		utils.HandleHTTPError(ctx, utils.Error(messages.FailedDataFromBody, http.StatusBadRequest, err))
		return
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Validation failed : %+v", err)
		message := utils.TranslateValidationError(c.Validate, err)
		utils.HandleHTTPError(ctx, utils.Error(message, http.StatusBadRequest, err))
		return
	}

	response, err := c.UseCase.Create(ctx.Request.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to create product : %+v", err)
		utils.HandleHTTPError(ctx, err)
		return
	}

	res := utils.SuccessResponse(messages.ProductCreated, response)
	ctx.JSON(http.StatusCreated, res)
}
