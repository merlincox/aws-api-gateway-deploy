package front

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/merlincox/aws-api-gateway-deploy/pkg/models"
)

func makeFront() Front {
	return NewFront(models.Status{}, 123)
}

func TestUnknownRoute(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	testFront := makeFront()

	Convey("When sending an request with an unknown path", t, func() {

		request := events.APIGatewayProxyRequest{
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: `/unknownpath`,
				HTTPMethod:   `GET`,
			},
		}

		Convey("Then it should return a bad request status code", func() {
			response, err := testFront.Handler(request)
			So(response.Body, ShouldEqual, `{"message":"No such route as GET/unknownpath","code":400}`)
			So(response.Headers["Access-Control-Allow-Origin"], ShouldEqual, "*")
			So(response.StatusCode, ShouldEqual, 400)
			So(err, ShouldBeNil)
		})
	})
}

func (front *Front) dummyDataRouter(route string) innerHandler {

	return front.dummyDataHandler
}

func (front Front) dummyDataHandler(request events.APIGatewayProxyRequest) (result interface{}, apiError models.ApiError) {

	return struct{Data string `json:"data"`}{Data: "Dummy"}, nil
}

func TestFrontDummyData(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	testFront := makeFront()

	testFront.router = testFront.dummyDataRouter

	Convey("When a handler returns data with no error", t, func() {

		request := events.APIGatewayProxyRequest{
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: `/whatever`,
				HTTPMethod:   `GET`,
			},
		}

		Convey("Then front should return a 200 request status code, and a JSON encoded string of the data", func() {
			response, err := testFront.Handler(request)
			So(response.Body, ShouldEqual, `{"data":"Dummy"}`)
			So(response.Headers["Access-Control-Allow-Origin"], ShouldEqual, "*")
			So(response.Headers["Cache-Control"], ShouldEqual, "max-age=123")
			So(response.StatusCode, ShouldEqual, 200)
			So(err, ShouldBeNil)
		})
	})
}

func (front *Front) errorRouter(route string) innerHandler {

	return front.errorHandler
}

func (front Front) errorHandler(request events.APIGatewayProxyRequest) (result interface{}, apiError models.ApiError) {

	return nil, models.ConstructApiError(345, "A simulated error: %v", "error")
}

func TestFrontErrorData(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	testFront := makeFront()

	testFront.router = testFront.errorRouter

	Convey("When a handler returns a non-nil ApiError", t, func() {

		request := events.APIGatewayProxyRequest{
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: `/whatever`,
				HTTPMethod:   `GET`,
			},
		}

		Convey("Then front should return the ApiError code and a JSON encoded error body with the ApiError message", func() {
			response, err := testFront.Handler(request)
			So(response.Body, ShouldEqual, `{"message":"A simulated error: error","code":345}`)
			So(response.Headers["Access-Control-Allow-Origin"], ShouldEqual, "*")
			So(response.StatusCode, ShouldEqual, 345)
			So(err, ShouldBeNil)
		})
	})
}

func (front *Front) unmarshallableRouter(route string) innerHandler {

	return front.unmarshallableHandler
}

func (front Front) unmarshallableHandler(request events.APIGatewayProxyRequest) (result interface{}, apiError models.ApiError) {

	return func(){}, nil
}

func TestFrontUnparseableData(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	testFront := makeFront()

	testFront.router = testFront.unmarshallableRouter

	Convey("When a handler returns unmarshallable data", t, func() {

		request := events.APIGatewayProxyRequest{
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: `/whatever`,
				HTTPMethod:   `GET`,
			},
		}

		Convey("Then front should return 500 and a JSON encoded error body with an 'Unmarshallable data' message", func() {
			response, err := testFront.Handler(request)
			So(response.Body, ShouldEqual, `{"message":"Unmarshallable data","code":500}`)
			So(response.Headers["Access-Control-Allow-Origin"], ShouldEqual, "*")
			So(response.StatusCode, ShouldEqual, 500)
			So(err, ShouldBeNil)
		})
	})
}

func (front *Front) panickyRouter(route string) innerHandler {

	return front.panickyHandler
}

func (front Front) panickyHandler(request events.APIGatewayProxyRequest) (result interface{}, apiError models.ApiError) {

	panic("Simulated panic")
	return
}

func TestFrontPanicRecovery(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	testFront := makeFront()

	testFront.router = testFront.panickyRouter

	Convey("When encountering a handler panic", t, func() {

		request := events.APIGatewayProxyRequest{
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: `/whatever`,
				HTTPMethod:   `GET`,
			},
		}

		Convey("Then front should return a 500 request status code and a JSON encoded error body with the panic message", func() {
			response, err := testFront.Handler(request)
			So(response.Body, ShouldEqual, `{"message":"Simulated panic","code":500}`)
			So(response.Headers["Access-Control-Allow-Origin"], ShouldEqual, "*")
			So(response.StatusCode, ShouldEqual, 500)
			So(err, ShouldBeNil)
		})
	})
}
