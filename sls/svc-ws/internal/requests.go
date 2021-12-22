package ws

import "encoding/json"

// WebSocketRequest is used for all custom routes, and allows switching
// on the specified action before unmarshalling the full request payload.
type WebSocketRequest struct {
	Action string          `json:"action"`
	Body   json.RawMessage `json:"body"`
}

type RequestSendMessage struct {
	SenderID   string `json:"senderId"`
	ReceiverID string `json:"receiverId"`
	Message    string `json:"message"`
}

type RequestJoin struct {
	Username string `json:"username"`
}
