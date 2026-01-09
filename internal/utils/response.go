package utils

import (
	"errors"
	"net/http"
	"snack-store-api/internal/messages"
	"snack-store-api/internal/model"

	"github.com/gin-gonic/gin"
)

func FailedResponse(code string, message string) model.WebResponse[any] {
	return model.WebResponse[any]{
		Error: &model.ErrorResponse{
			Code:    code,
			Message: message,
		},
	}
}

type HTTPError interface {
	error
	Status() int
	Message() string
	Unwrap() error
}

type appError struct {
	msg    string
	status int
	err    error
}

func (e *appError) Error() string   { return e.msg }
func (e *appError) Message() string { return e.msg }
func (e *appError) Status() int     { return e.status }
func (e *appError) Unwrap() error   { return e.err }

func Error(message string, status int, err error) error {
	return &appError{msg: message, status: status, err: err}
}

func HandleHTTPError(ctx *gin.Context, err error) {
	var httpErr HTTPError

	if errors.As(err, &httpErr) {
		code := errorCodeFromStatus(httpErr.Status())
		res := FailedResponse(code, httpErr.Message())
		ctx.AbortWithStatusJSON(httpErr.Status(), res)
		return
	}

	res := FailedResponse(errorCodeFromStatus(http.StatusInternalServerError), messages.InternalServerError)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
}

func SuccessResponse[T any](message string, data T) model.WebResponse[T] {
	return model.WebResponse[T]{
		Data:    data,
		Message: message,
	}
}

func SuccessWithPaginationResponse[T any](
	message string,
	data []T,
	paging model.PageMetadata,
) model.WebResponse[[]T] {
	return model.WebResponse[[]T]{
		Message: message,
		Data:    data,
		Paging:  &paging,
	}
}

func errorCodeFromStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "VALIDATION_ERROR"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusTooManyRequests:
		return "TOO_MANY_REQUESTS"
	case http.StatusUnprocessableEntity:
		return "UNPROCESSABLE_ENTITY"
	default:
		return "INTERNAL_SERVER_ERROR"
	}
}
