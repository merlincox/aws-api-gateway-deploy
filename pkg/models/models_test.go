package models

import (
	"errors"
	"testing"

	"github.com/merlincox/aws-api-gateway-deploy/pkg/utils"
)

func TestApiErrorBody(t *testing.T) {

	eb := ApiErrorBody{
		Message: "Testing testing 1 2 3",
		Code:    123,
	}

	tb := errBody{
		body: eb,
	}

	utils.AssertEquals(t, "ApiErrorBody error", "Testing testing 1 2 3", tb.Error())
	utils.AssertEquals(t, "ApiErrorBody code", 123, tb.StatusCode())
	utils.AssertEquals(t, "ApiErrorBody body", eb, tb.ErrorBody())
}

func TestErrorf(t *testing.T) {

	var err ApiError

	err = ConstructApiError(123, "Testing testing %v %v %v", 1, "2", 3)

	eb := ApiErrorBody{
		Message: "Testing testing 1 2 3",
		Code:    123,
	}

	utils.AssertEquals(t, "ConstructApiError string", "Testing testing 1 2 3", err.Error())
	utils.AssertEquals(t, "ConstructApiError code", 123, err.StatusCode())
	utils.AssertEquals(t, "ConstructApiError body", eb, err.ErrorBody())
}

func TestErrorWrap(t *testing.T) {

	innerErr := errors.New("I am an error")

	innerErr2 := ConstructApiError(123, "I am an %v", "API error")

	err := ErrorWrap(innerErr)

	err2 := ErrorWrap(innerErr2)

	errBody := ApiErrorBody{
		Message: "I am an error",
		Code:    500,
	}

	errBody2 := ApiErrorBody{
		Message: "I am an API error",
		Code:    123,
	}

	utils.AssertEquals(t, "Non API error string", "I am an error", err.Error())
	utils.AssertEquals(t, "Non API error code", 500, err.StatusCode())
	utils.AssertEquals(t, "Non API error body", errBody, err.ErrorBody())
	utils.AssertEquals(t, "API error string", "I am an API error", err2.Error())
	utils.AssertEquals(t, "API error code", 123, err2.StatusCode())
	utils.AssertEquals(t, "API error body", errBody2, err2.ErrorBody())
}
