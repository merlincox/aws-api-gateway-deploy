package front

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/merlincox/aws-api-gateway-deploy/pkg/models"
	"github.com/merlincox/aws-api-gateway-deploy/pkg/utils"
	"fmt"
)

var testFront = NewFront(models.Status{
	Branch:    "testing",
	Platform:  "test",
	Commit:    "a00eaaf45694163c9b728a7b5668e3d510eb3eb0",
	Release:   "v1.0.1",
	Timestamp: "2019-01-02T14:52:36.951375973Z",
}, 123)


func TestStatusRoute(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	expected := models.Status{
		Branch:    "testing",
		Platform:  "test",
		Commit:    "a00eaaf45694163c9b728a7b5668e3d510eb3eb0",
		Release:   "v1.0.1",
		Timestamp: "2019-01-02T14:52:36.951375973Z",
	}

	Convey("When sending an request with the /status route", t, func() {

		request := events.APIGatewayProxyRequest{
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: `/status`,
				HTTPMethod:   `GET`,
			},
		}

		Convey("Then it should return the status", func() {
			response, err := testFront.Handler(request)
			So(response.Body, ShouldEqual, utils.JsonStringify(expected))
			So(response.Headers["Access-Control-Allow-Origin"], ShouldEqual, "*")
			So(response.Headers["Cache-Control"], ShouldEqual, "max-age=123")
			So(response.StatusCode, ShouldEqual, 200)
			So(err, ShouldBeNil)
		})
	})
}

func testCalc(t *testing.T, val1, val2, result float64, op, fullop string) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	expected := models.CalculationResult{
		Locale: "en-GB",
		Op:     fullop,
		Result: result,
		Val1:   val1,
		Val2:   val2,
	}

	Convey(fmt.Sprintf("When sending an request with the /calc route with %v operator", fullop), t, func() {

		request := events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"val1": fmt.Sprintf("%v", val1),
				"val2": fmt.Sprintf("%v", val2),
			},
			PathParameters: map[string]string{
				"op": op,
			},
			Headers: map[string]string{
				"x-locale": "en-GB",
			},
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: `/calc/{op}`,
				HTTPMethod:   `GET`,
			},
		}

		Convey("Then it should return the correct result", func() {
			response, err := testFront.Handler(request)
			So(response.Body, ShouldEqual, utils.JsonStringify(expected))
			So(response.Headers["Access-Control-Allow-Origin"], ShouldEqual, "*")
			So(response.Headers["Cache-Control"], ShouldEqual, "max-age=123")
			So(response.StatusCode, ShouldEqual, 200)
			So(err, ShouldBeNil)
		})
	})

}

func TestCalcRouteAdd(t *testing.T) {
	testCalc(t, 3.5, 2.25, 5.75, "add", "add")
}

func TestCalcRouteSub(t *testing.T) {
	testCalc(t, 3.5, 2.25, 1.25, "sub", "subtract")
}


func TestCalcRouteMult(t *testing.T) {
	testCalc(t, 1.5, 7, 10.5, "mul", "multiply")
}

func TestCalcRoutePow(t *testing.T) {
	testCalc(t, 2, 3, 8, "pow", "power")
}

func TestCalcRouteRoot(t *testing.T) {

	testCalc(t, 16, 2, 4, "roo", "root")
}

func testCalcRouteBad(t *testing.T, val1, val2 float64, op, context, msg string) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	expected := models.ApiErrorBody{
		Message: msg,
		Code:    400,
	}

	Convey(context, t, func() {

		request := events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"val1": fmt.Sprintf("%v", val1),
				"val2": fmt.Sprintf("%v", val2),
			},
			PathParameters: map[string]string{
				"op": op,
			},
			Headers: map[string]string{
				"x-locale": "en-GB",
			},
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: `/calc/{op}`,
				HTTPMethod:   `GET`,
			},
		}

		Convey("Then it should return correct error", func() {
			response, err := testFront.Handler(request)
			So(response.Body, ShouldEqual, utils.JsonStringify(expected))
			So(response.Headers["Access-Control-Allow-Origin"], ShouldEqual, "*")
			So(response.Headers["Cache-Control"], ShouldEqual, "max-age=123")
			So(response.StatusCode, ShouldEqual, 400)
			So(err, ShouldBeNil)
		})
	})
}

func TestCalcRouteBadOp(t *testing.T) {
	testCalcRouteBad(t, 1,2, "bad", "When sending a request to the /calc route with a bad operator", "Unknown calc operation: bad")
}

func TestCalcRouteInf(t *testing.T) {
	testCalcRouteBad(t, 1,0, "div", "When sending a request to the /calc route with inf result", "Out of limits: 1 divide 0")
}

func TestCalcRouteNegInf(t *testing.T) {
	testCalcRouteBad(t, -1,0, "div", "When sending a request to the /calc route with negative inf result", "Out of limits: -1 divide 0")
}

func TestCalcRouteNaN(t *testing.T) {
	testCalcRouteBad(t, -1,2, "root", "When sending a request to the /calc route with NaN result", "Out of limits: -1 root 2")
}