package front

import (
	"github.com/aws/aws-lambda-go/events"

	"github.com/merlincox/aws-api-gateway-deploy/pkg/models"
	"strconv"
	"fmt"
	"math"
)

func (front Front) statusHandler(request events.APIGatewayProxyRequest) (interface{}, models.ApiError) {

	return front.status, nil
}

func (front Front) calcHandler(request events.APIGatewayProxyRequest) (interface{}, models.ApiError) {

	op := request.PathParameters["op"]

	val1, err := getFloatFromRequest(request, "val1")

	if err != nil {
		return nil, models.ConstructApiError(400, err.Error())
	}

	val2, err := getFloatFromRequest(request, "val2")

	if err != nil {
		return nil, models.ConstructApiError(400, err.Error())
	}

	locale, ok := request.Headers["locale"]

	if !ok {
		locale = "undefined"
	}

	result := models.CalculationResult{
		Locale: locale,
		Op:     op,
		Val1:   val1,
		Val2:   val2,
	}

	switch op[0:3] {

	case "add":

		result.Result = val1 + val2

	case "sub":

		result.Result = val1 - val2

	case "mul":

		result.Result = val1 * val2

	case "div":

		result.Result = val1 / val2

	case "pow":

		result.Result = math.Pow(val1, val2)

	case "roo":

		result.Result = math.Pow(val1, 1/val2)

	default:

		return nil, models.ConstructApiError(400, "Unknown operation %v", op)
	}

	return result, nil
}

func getFloatFromRequest(request events.APIGatewayProxyRequest, key string) (result float64, err error) {

	val, ok := request.QueryStringParameters[key]

	if ! ok {
		err = fmt.Errorf("Missing parameter %v", key)
		return
	}

	result, err = strconv.ParseFloat(val, 64)

	return
}