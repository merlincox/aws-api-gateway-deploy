package front

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/merlincox/aws-api-gateway-deploy/pkg/models"
	"github.com/merlincox/aws-api-gateway-deploy/pkg/utils"
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

func TestCalcRouteAdd(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	expected := models.CalculationResult{
		Locale: "en-GB",
		Op:     "add",
		Result: 3.75,
		Val1:   1.5,
		Val2:   2.25,
	}

	Convey("When sending an request with the /calc route", t, func() {

		request := events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"val1": "1.5",
				"val2": "2.25",
			},
			PathParameters: map[string]string{
				"op": "add",
			},
			Headers: map[string]string{
				"x-locale": "en-GB",
			},
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: `/calc/{op}`,
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

func TestCalcRouteSub(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	expected := models.CalculationResult{
		Locale: "en-GB",
		Op:     "subtract",
		Result: 1.25,
		Val1:   3.5,
		Val2:   2.25,
	}

	Convey("When sending an request with the /calc route", t, func() {

		request := events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"val1": "3.5",
				"val2": "2.25",
			},
			PathParameters: map[string]string{
				"op": "subtract",
			},
			Headers: map[string]string{
				"x-locale": "en-GB",
			},
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: `/calc/{op}`,
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


func TestCalcRouteMult(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	expected := models.CalculationResult{
		Locale: "en-GB",
		Op:     "multiply",
		Result: 10.5,
		Val1:   1.5,
		Val2:   7,
	}

	Convey("When sending an request with the /calc route", t, func() {

		request := events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"val1": "1.5",
				"val2": "7",
			},
			PathParameters: map[string]string{
				"op": "multiply",
			},
			Headers: map[string]string{
				"x-locale": "en-GB",
			},
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: `/calc/{op}`,
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

func TestCalcRoutePow(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	expected := models.CalculationResult{
		Locale: "en-GB",
		Op:     "power",
		Result: 8,
		Val1:   2,
		Val2:   3,
	}

	Convey("When sending an request with the /calc route", t, func() {

		request := events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"val1": "2",
				"val2": "3",
			},
			PathParameters: map[string]string{
				"op": "power",
			},
			Headers: map[string]string{
				"x-locale": "en-GB",
			},
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: `/calc/{op}`,
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

func TestCalcRouteRoot(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	expected := models.CalculationResult{
		Locale: "en-GB",
		Op:     "root",
		Result: 4,
		Val1:   16,
		Val2:   2,
	}

	Convey("When sending an request with the /calc route", t, func() {

		request := events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"val1": "16",
				"val2": "2",
			},
			PathParameters: map[string]string{
				"op": "root",
			},
			Headers: map[string]string{
				"x-locale": "en-GB",
			},
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: `/calc/{op}`,
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

func TestCalcRouteBadOp(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	expected := models.ApiErrorBody{
		Message: "Unknown calc operation: bad",
		Code:    400,
	}

	Convey("When sending an request with the /calc route", t, func() {

		request := events.APIGatewayProxyRequest{
			QueryStringParameters: map[string]string{
				"val1": "1.5",
				"val2": "2.25",
			},
			PathParameters: map[string]string{
				"op": "bad",
			},
			Headers: map[string]string{
				"x-locale": "en-GB",
			},
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: `/calc/{op}`,
				HTTPMethod:   `GET`,
			},
		}

		Convey("Then it should return the status", func() {
			response, err := testFront.Handler(request)
			So(response.Body, ShouldEqual, utils.JsonStringify(expected))
			So(response.Headers["Access-Control-Allow-Origin"], ShouldEqual, "*")
			So(response.Headers["Cache-Control"], ShouldEqual, "max-age=123")
			So(response.StatusCode, ShouldEqual, 400)
			So(err, ShouldBeNil)
		})
	})
}