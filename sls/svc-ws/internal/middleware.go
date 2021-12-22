package ws

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type HandlerFunc func(ctx context.Context, req *events.APIGatewayWebsocketProxyRequest) (Response, error)

// Middleware executed initialization logic common to all WS handlers.
func Middleware(logger *Logger, next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, req *events.APIGatewayWebsocketProxyRequest) (Response, error) {
		// set common context values for logging
		ctx = context.WithValue(ctx, KeyConnectionID, req.RequestContext.ConnectionID)
		ctx = context.WithValue(ctx, KeyRequestID, req.RequestContext.RequestID)

		// flush buffered logs on exit
		defer logger.Sync()
		logger.Info(ctx, "request to "+req.RequestContext.RouteKey)

		return next(ctx, req)
	}
}
