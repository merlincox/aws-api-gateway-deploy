package models

import (
	"fmt"
	"net/http"
)

type ApiError interface {
	Error() string
	StatusCode() int
	ErrorBody() ApiErrorBody
}

type ApiErrorBody struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type errBody struct {
	body ApiErrorBody
}

func (err errBody) Error() string {
	return err.body.Message
}

func (err errBody) StatusCode() int {
	return err.body.Code
}

func (err errBody) ErrorBody() ApiErrorBody {
	return err.body
}

func ConstructApiError(code int, format string, a ...interface{}) ApiError {

	return errBody{
		body: ApiErrorBody{
			Message: fmt.Sprintf(format, a...),
			Code:    code,
		},
	}
}

func ErrorWrap(err error) ApiError {

	apiErr, ok := err.(ApiError)

	if ok {
		return apiErr
	}

	return errBody{
		body: ApiErrorBody{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		},
	}
}
