package utils

import (
	"errors"
	"net/http"
	"snack-store-api/internal/messages"
	"snack-store-api/internal/model"

	"github.com/gin-gonic/gin"
)

func FailedResponse(message string) model.WebResponse[any] {
	return model.WebResponse[any]{
		Errors: message,
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
		res := FailedResponse(httpErr.Message())
		ctx.AbortWithStatusJSON(httpErr.Status(), res)
		return
	}

	res := FailedResponse(messages.InternalServerError)
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
