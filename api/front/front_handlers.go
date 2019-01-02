package front

import (
	"github.com/aws/aws-lambda-go/events"

	"github.com/merlincox/aws-api-gateway-deploy/pkg/models"
)

func (front Front) statusHandler(request events.APIGatewayProxyRequest) (interface{}, models.ApiError) {

	return front.status, nil
}

