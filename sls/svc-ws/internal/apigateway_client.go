package ws

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"

	"github.com/Jimeux/sls-chat-app/sls/lib/awsiface"
)

// APIClient is a simple wrapper around APIGatewayManagementAPI.
type APIClient struct {
	client awsiface.APIGatewayManagementAPI
}

func NewAPIGatewayClient(client awsiface.APIGatewayManagementAPI) *APIClient {
	return &APIClient{client: client}
}

// NewAPIClientFromConfig create a APIClient instance from a given aws.Config instance.
// stage and domain are used to construct the required endpoint resolver.
func NewAPIClientFromConfig(cfg aws.Config, stage, domainName string) *APIClient {
	var endpoint url.URL
	endpoint.Scheme = "https"
	endpoint.Path = stage
	endpoint.Host = domainName
	endpointResolver := apigatewaymanagementapi.EndpointResolverFromURL(endpoint.String())
	return NewAPIGatewayClient(apigatewaymanagementapi.NewFromConfig(
		cfg, apigatewaymanagementapi.WithEndpointResolver(endpointResolver)))
}

// Publish posts data to the connection identified by connID.
func (c *APIClient) Publish(ctx context.Context, connID string, data []byte) error {
	_, err := c.client.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
		Data:         data,
		ConnectionId: aws.String(connID),
	})
	return err
}
