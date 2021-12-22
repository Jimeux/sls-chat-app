package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi/types"
)

// WebSocketService implements the core logic of svc-ws.
type WebSocketService struct {
	client     *Client
	repository *Repository
}

func NewWebSocketService(client *Client, repository *Repository) *WebSocketService {
	return &WebSocketService{
		client:     client,
		repository: repository,
	}
}

// Connect saves a reference to connID in the datastore.
func (s *WebSocketService) Connect(ctx context.Context, connID string) (Response, error) {
	// save connection initially without username
	if err := s.repository.SaveConnection(ctx, connID, ""); err != nil {
		return ResponseInternalServerError(), fmt.Errorf("failed to connect connectionID %s: %w", connID, err)
	}
	return ResponseOK(), nil
}

func (s *WebSocketService) Disconnect(ctx context.Context, connID string) (Response, error) {
	if err := s.repository.RemoveConnection(ctx, connID); err != nil {
		return ResponseInternalServerError(), err
	}

	connections, err := s.repository.GetConnections(ctx)
	if err != nil {
		return ResponseInternalServerError(), err
	}

	res := WebSocketPayload{
		Action: ActionUserLeft,
		Body:   PayloadUserLeft{ConnectionID: connID},
	}
	b, err := json.Marshal(&res)
	if err != nil {
		return ResponseInternalServerError(), err
	}

	for _, conn := range connections {
		// skip current user in case of read consistency latency
		if conn.ConnectionID == connID {
			continue
		}
		if err := s.publish(ctx, b, conn.ConnectionID); err != nil {
			return ResponseInternalServerError(), err
		}
	}
	return ResponseOK(), nil
}

func (s *WebSocketService) Message(ctx context.Context, connID string, req RequestSendMessage) (Response, error) {
	// validate request
	if req.ReceiverID == "" {
		return ResponseBadRequest("Invalid channelId"), nil
	}
	if req.Message == "" {
		return ResponseBadRequest("Invalid message"), nil
	}

	// fetch recipients connection ID
	conn, err := s.repository.GetConnection(ctx, req.ReceiverID)
	if err != nil {
		return ResponseInternalServerError(), err
	}

	res := WebSocketPayload{
		Action: ActionUserMessage,
		Body: PayloadUserMessage{
			SenderID:   connID,
			ReceiverID: req.ReceiverID,
			Message:    req.Message,
		},
	}
	b, err := json.Marshal(&res)
	if err != nil {
		return ResponseInternalServerError(), err
	}

	// publish to both recipient and sender
	if err := s.publish(ctx, b, conn.ConnectionID, connID); err != nil {
		return ResponseInternalServerError(), err
	}
	return ResponseOK(), nil
}

func (s *WebSocketService) Join(ctx context.Context, connID string, req RequestJoin) (Response, error) {
	if req.Username == "" {
		return ResponseBadRequest("Invalid username"), nil
	}

	if err := s.repository.SaveConnection(ctx, connID, req.Username); err != nil {
		return ResponseInternalServerError(), fmt.Errorf("failed to save connection for join: %w", err)
	}

	users, err := s.repository.GetConnections(ctx)
	if err != nil {
		return ResponseInternalServerError(), err
	}

	body := PayloadUserJoined{UserList: users}
	res := WebSocketPayload{
		Action: ActionUserJoined,
	}

	// send user list to all connected users
	for _, u := range users {
		body.ConnectionId = u.ConnectionID
		res.Body = body
		b, err := json.Marshal(&res)
		if err != nil {
			return ResponseInternalServerError(), fmt.Errorf("failed to marshal user list: %w", err)
		}

		if err := s.publish(ctx, b, u.ConnectionID); err != nil {
			return ResponseInternalServerError(), fmt.Errorf("failed to join: %w", err)
		}
	}
	return ResponseOK(), nil
}

// publish posts the given data to one or more connections. When encountering
// GoneException (a stale connection ID), the corresponding connection is
// automatically removed from the database and the error is ignored.
func (s *WebSocketService) publish(ctx context.Context, data []byte, connIDs ...string) error {
	for _, connID := range connIDs {
		if err := s.client.Publish(ctx, connID, data); err != nil {
			var errGone *types.GoneException
			if errors.As(err, &errGone) {
				if err := s.repository.RemoveConnection(ctx, connID); err != nil {
					return fmt.Errorf("failed to remove connection connID %s: %w", connID, err)
				}
				continue
			}
			return fmt.Errorf("failed to post to connection connID %s: %w", connID, err)
		}
	}
	return nil
}
