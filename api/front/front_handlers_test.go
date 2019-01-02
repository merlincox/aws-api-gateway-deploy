package front

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/merlincox/aws-api-gateway-deploy/pkg/models"
	"github.com/merlincox/aws-api-gateway-deploy/pkg/utils"
)

func makeStatusFront(status models.Status) Front {
	return NewFront(status, 123)
}

func TestStatusRoute(t *testing.T) {

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	status := models.Status{
		Branch:    "testing",
		Platform:  "test",
		Commit:    "a00eaaf45694163c9b728a7b5668e3d510eb3eb0",
		Release:   "v1.0.1",
		Timestamp: "2019-01-02T14:52:36.951375973Z",
	}

	testFront := makeStatusFront(status)

	Convey("When sending an request with an unknown path", t, func() {

		request := events.APIGatewayProxyRequest{
			RequestContext: events.APIGatewayProxyRequestContext{
				ResourcePath: `/status`,
				HTTPMethod:   `GET`,
			},
		}

		Convey("Then it should return a bad request status code", func() {
			response, err := testFront.Handler(request)
			So(response.Body, ShouldEqual, utils.JsonStringify(status))
			So(response.Headers["Access-Control-Allow-Origin"], ShouldEqual, "*")
			So(response.StatusCode, ShouldEqual, 200)
			So(err, ShouldBeNil)
		})
	})
}


