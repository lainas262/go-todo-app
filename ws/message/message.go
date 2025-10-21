package message

import (
	"context"
	"encoding/json"

	"github.com/coder/websocket"
)

type Client interface {
	Send(message []byte)
	Ctx() context.Context
	Conn() *websocket.Conn
}

type Room interface {
	Collaborators() map[Client]bool
	Broadcast(Event)
}

type Event struct {
	Name    string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
	Client  Client
}

type OpenTodoPayload struct {
	TodoId int64 `json:"todo_id"`
}

type UpdateTodoPayload struct {
	TodoId int64  `json:"todo_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
