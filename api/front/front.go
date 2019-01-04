// The Front package routes HTTP requests to an appropriately routed handler bound to ReelIndexer and JsonHydrator instances and returns a response
// whose Body member is a JSON-encoded API object. In case of error it will be a JSON-encoded ApiErrorBody.
package front

import (
	"fmt"
	"net/http"
	"time"
	"log"
	"strconv"
	"runtime/debug"

	"github.com/aws/aws-lambda-go/events"

	"github.com/merlincox/aws-api-gateway-deploy/pkg/utils"
	"github.com/merlincox/aws-api-gateway-deploy/pkg/models"
)

type Front struct {
	status      models.Status
	router      func(route string) innerHandler
	cacheMaxAge int
}


type FrontHandler func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
type innerHandler func(request events.APIGatewayProxyRequest) (interface{}, models.ApiError)

// NewFront Create a new Front object
//
func NewFront(status models.Status, cacheMaxAge int) Front {

	f := Front{
		status: status,
		cacheMaxAge: cacheMaxAge,
	}

	f.router = f.getHandlerForRoute

	return f
}

// Receive a APIGatewayProxyRequest and returns a APIGatewayProxyResponse with nil error
//
// Any panic should be recovered and wrapped into an ApiErrorBody, and the trace logged
func (front Front) Handler(request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {

	defer func() {

		if r := recover(); r != nil {
			log.Println(utils.JsonStack(r, debug.Stack()))
			apiErr := models.ConstructApiError(http.StatusInternalServerError, "%v", r)
			response = front.buildResponse(nil, apiErr)
		}

	}()

	route := getRoute(request)
	log.Println("Handling a request for %v.", route)

	response = front.buildResponse(front.router(route)(request))

	return
}

func (front *Front) getHandlerForRoute(route string) innerHandler {

	switch route {

	case "GET/status":
		return front.statusHandler

	case "GET/calc/{op}":
		return front.calcHandler

	}

	return front.unknownRouteHandler
}

func getRoute(request events.APIGatewayProxyRequest) string {

	return request.RequestContext.HTTPMethod + request.RequestContext.ResourcePath
}

func (front Front) unknownRouteHandler(request events.APIGatewayProxyRequest) (interface{}, models.ApiError) {

	return nil, models.ConstructApiError(http.StatusNotFound, "No such route as %v", getRoute(request))
}

func (front *Front) seoString(s string) string {
	return utils.Slug(s)
}

func (front *Front) buildResponse(data interface{}, err models.ApiError) events.APIGatewayProxyResponse {

	var (
		body       string
		statusCode int
	)

	if err != nil {

		body = utils.JsonStringify(err.ErrorBody())
		statusCode = err.StatusCode()
		log.Printf("ERROR: Returning %v: %v", statusCode, err.Error())

	} else {

		body = utils.JsonStringify(data)
		statusCode = http.StatusOK
	}

	// handle unlikely case where json.Marshall fails for the data argument
	if body == "" {
		statusCode = http.StatusInternalServerError
		body = fmt.Sprintf(`{"message":"Unmarshallable data","code":%v}`, statusCode)
		log.Printf("ERROR: Returning %v: %v", statusCode, "Unmarshallable data")
	}

	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: statusCode,
		Headers: map[string]string{
			"Cache-Control":               "max-age=" + strconv.Itoa(front.cacheMaxAge),
			"Access-Control-Allow-Origin": "*",
			"X-Timestamp":        time.Now().UTC().Format(time.RFC3339Nano),
		}}
}
