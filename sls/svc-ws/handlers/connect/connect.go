package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"

	ws "github.com/Jimeux/sls-chat-app/sls/svc-ws/internal"
)

// Package-level vars survive for the life of the
// Lambda, and help avoid unnecessary allocations.
var (
	logger *ws.Logger
	svc    *ws.WebSocketService
)

// main initialises package-level vars and calls lambda.Start, passing
// handler, which is wrapped in middleware that initialises logging.
func main() {
	logger = ws.NewLogger()
	cf := ws.NewConfig()
	cfg, _ := config.LoadDefaultConfig(context.Background())
	repository := ws.NewRepositoryFromConfig(cfg, cf.ConnectionsTable)
	client := ws.NewAPIClientFromConfig(cfg, cf.Stage, cf.APIGatewayDomain)
	svc = ws.NewWebSocketService(client, repository)

	lambda.Start(ws.Middleware(logger, handler))
}

// handler delegates to WebSocketService to avoid implementing core logic itself.
func handler(ctx context.Context, event *events.APIGatewayWebsocketProxyRequest) (ws.Response, error) {
	res, err := svc.Connect(ctx, event.RequestContext.ConnectionID)
	if err != nil {
		logger.Error(ctx, "connect failure", err)
		return res, err
	}
	logger.Info(ctx, "connect success")
	return res, nil
}
