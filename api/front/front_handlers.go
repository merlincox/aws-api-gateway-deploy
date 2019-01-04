package front

import (
	"strconv"
	"fmt"
	"math"

	"golang.org/x/text/message"
	"golang.org/x/text/language"

	"github.com/aws/aws-lambda-go/events"

	"github.com/merlincox/aws-api-gateway-deploy/pkg/models"
)

func (front Front) statusHandler(request events.APIGatewayProxyRequest) (interface{}, models.ApiError) {

	return front.status, nil
}

func (front Front) calcHandler(request events.APIGatewayProxyRequest) (interface{}, models.ApiError) {

	var (
		result float64
		fullop string
	)

	locale, ok := request.Headers["Accept-Language"]

	p := message.NewPrinter(language.Make(locale))

	if !ok {
		locale = "undefined"
	}

	op := request.PathParameters["op"]

	val1, err := getFloatFromRequest(request, "val1")

	if err != nil {
		return nil, models.ConstructApiError(400, err.Error())
	}

	val2, err := getFloatFromRequest(request, "val2")

	if err != nil {
		return nil, models.ConstructApiError(400, err.Error())
	}

	switch op[0:3] {

	case "add":

		result = val1 + val2
		fullop = "add"

	case "sub":

		result = val1 - val2
		fullop = "subtract"

	case "mul":

		result = val1 * val2
		fullop = "multiply"

	case "div":

		result = val1 / val2
		fullop = "divide"

	case "pow":

		result = math.Pow(val1, val2)
		fullop = "power"

	case "roo":

		result = math.Pow(val1, 1/val2)
		fullop = "root"

	default:

		return nil, models.ConstructApiError(400, "Unknown calc operation: %v", op)
	}

	if math.IsNaN(result) || math.IsInf(result, 1) || math.IsInf(result, -1) {
		return nil, models.ConstructApiError(400, "Out of limits: %v %v %v", val1, fullop, val2)
	}

	return models.CalculationResult{
		Locale: locale,
		Op:     fullop,
		Val1:   val1,
		Val2:   val2,
		Result: p.Sprintf("%v", result),
	}, nil
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
