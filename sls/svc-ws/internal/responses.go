package ws

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

const (
	ActionError       = "error"
	ActionSuccess     = "success"
	ActionUserJoined  = "user:joined"
	ActionUserLeft    = "user:left"
	ActionUserMessage = "user:message"
)

type Response = events.APIGatewayProxyResponse

func ResponseInternalServerError() Response {
	return Response{
		StatusCode: http.StatusInternalServerError,
		Body:       `{"action": "` + ActionError + `", "code": 500, "message": "Unknown error"}`,
	}
}

func ResponseBadRequest(msg string) Response {
	return Response{
		StatusCode: http.StatusBadRequest,
		Body:       `{"action": "` + ActionError + `", "code": 400, "message": "` + msg + `"}`,
	}
}

func ResponseOK() Response {
	return Response{
		StatusCode: http.StatusOK,
		Body:       `{"action": "` + ActionSuccess + `", "code": 200, "message": "OK"}`,
	}
}

type WebSocketPayload struct {
	Action string      `json:"action"`
	Body   interface{} `json:"body"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PayloadUserJoined struct {
	ConnectionId string        `json:"connectionId"`
	UserList     []*Connection `json:"userList"`
}

type PayloadUserLeft struct {
	ConnectionID string `json:"connectionId"`
}

type PayloadUserMessage struct {
	SenderID   string `json:"senderId"`
	ReceiverID string `json:"receiverId"`
	Message    string `json:"message"`
}
