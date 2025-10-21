package handler

import (
	"context"
	"encoding/json"
	"example/todolist/ws/message"

	"log"
)

type Hub interface {
	GetRoom(id int64) message.Room
}

type OpenTodoHandler struct {
	Hub Hub
}

func (h *OpenTodoHandler) Handle(ctx context.Context, event message.Event) {
	log.Printf("Trigger Open Handler call")
	//register todo in redis as open with an exp time of 30 seconds for inactivity maybe?

	// var e Event = Event{
	// 	Name: "todo_opened",
	// 	Payload: ,
	// }
	// Client.Send <-

}

type UpdateTodoHandler struct {
	Hub Hub
}

func (h *UpdateTodoHandler) Handle(ctx context.Context, event message.Event) {
	log.Printf("Trigger Update Handler call")
	var payload message.UpdateTodoPayload
	err := json.Unmarshal(event.Payload, &payload)

	if err != nil {
		log.Println("Error parsing JSON payload - UpdateTodoHandler")
	}

	//avoid payload validation checks for now

	room := h.Hub.GetRoom(payload.TodoId)
	event.Name = "update_todo_broadcast"
	room.Broadcast(event)
}
