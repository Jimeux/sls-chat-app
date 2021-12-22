package awsiface

import (
	"context"

	apigw "github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

type APIGatewayManagementAPI interface {
	PostToConnection(
		ctx context.Context, params *apigw.PostToConnectionInput, optFns ...func(*apigw.Options),
	) (*apigw.PostToConnectionOutput, error)
}

type MockAPIGatewayManagementAPI struct {
	APIGatewayManagementAPI
	PostToConnectionFn func(
		ctx context.Context, params *apigw.PostToConnectionInput, optFns ...func(*apigw.Options),
	) (*apigw.PostToConnectionOutput, error)
	PostToConnectionInvoked bool
}

func (m *MockAPIGatewayManagementAPI) PostToConnection(
	ctx context.Context, params *apigw.PostToConnectionInput, optFns ...func(*apigw.Options)) (*apigw.PostToConnectionOutput, error) {
	m.PostToConnectionInvoked = true
	return m.PostToConnectionFn(ctx, params, optFns...)
}
