package hub

import (
	"context"
	"example/todolist/ws/message"
)

type EventHandler interface {
	Handle(ctx context.Context, event message.Event)
}
