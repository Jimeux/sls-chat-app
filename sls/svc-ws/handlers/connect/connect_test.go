package main

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/Jimeux/sls-chat-app/sls/lib/awsiface"
	ws "github.com/Jimeux/sls-chat-app/sls/svc-ws/internal"
)

func TestHandler(t *testing.T) {
	logger = ws.NewNopLogger()

	t.Run("returns 200 on successful save", func(t *testing.T) {
		ddb := &awsiface.MockDDB{}
		ddb.PutItemFn = func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		}
		repo := ws.NewRepository(ddb, "connections")
		svc = ws.NewWebSocketService(nil, repo)

		res, err := handler(context.Background(), &events.APIGatewayWebsocketProxyRequest{})
		if err != nil {
			t.Fatalf("unexpected error: %+v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("did not return valid response: got %d but expected %d", res.StatusCode, http.StatusOK)
		}
	})
	t.Run("returns 500 on save error", func(t *testing.T) {
		ddb := &awsiface.MockDDB{}
		ddb.PutItemFn = func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return nil, errors.New("save error")
		}
		repo := ws.NewRepository(ddb, "connections")
		svc = ws.NewWebSocketService(nil, repo)

		res, err := handler(context.Background(), &events.APIGatewayWebsocketProxyRequest{})
		if err == nil {
			t.Fatalf("expected error but got none")
		}
		if res.StatusCode != http.StatusInternalServerError {
			t.Fatalf("did not return valid response: got %d but expected %d", res.StatusCode, http.StatusInternalServerError)
		}
	})
}
