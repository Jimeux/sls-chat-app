package ws

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"

	"github.com/Jimeux/sls-chat-app/sls/lib/awsiface"
)

// Client is a simple wrapper around APIGatewayManagementAPI.
type Client struct {
	client awsiface.APIGatewayManagementAPI
}

func NewAPIGatewayClient(client awsiface.APIGatewayManagementAPI) *Client {
	return &Client{client: client}
}

// NewAPIGatewayClientFromConfig create a Client instance from a given aws.Config instance.
// The stage and domain name values are used to construct the endpoint resolver required by
// apigatewaymanagementapi.Client.
func NewAPIGatewayClientFromConfig(cfg aws.Config, stage, domainName string) *Client {
	var endpoint url.URL
	endpoint.Scheme = "https"
	endpoint.Path = stage
	endpoint.Host = domainName
	endpointResolver := apigatewaymanagementapi.EndpointResolverFromURL(endpoint.String())

	return NewAPIGatewayClient(apigatewaymanagementapi.NewFromConfig(
		cfg, apigatewaymanagementapi.WithEndpointResolver(endpointResolver)))
}

// Publish posts data to the WebSocket connection identified by connID.
func (c *Client) Publish(ctx context.Context, connID string, data []byte) error {
	_, err := c.client.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
		Data:         data,
		ConnectionId: aws.String(connID),
	})
	return err
}
