package handler

import (
	"example/todolist/middleware"
	"example/todolist/response"
	"example/todolist/ws/hub"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type WebsocketHandler struct {
	hub *hub.Hub
}

func CreateWebsocketHandler(service *hub.Hub) *WebsocketHandler {
	return &WebsocketHandler{hub: service}
}

func (h *WebsocketHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIdKey).(int64)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	log.Printf("Websocket Subscribe - Request received: %v", userId)

	token := r.FormValue("token")

	wsConn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"http://localhost:3000"},
	})

	if err != nil {
		log.Printf("Websocket accept err: %v", err)
		return
	}

	// passing token to hit api later. Ideally this would be handle with a gateway layer for auth
	// and internal service communication wouldn't need this perhaps
	client := hub.NewClient(token, userId, wsConn)

	h.hub.RegisterClient(client)
	go h.hub.CreateClientRooms(client)
	go h.hub.ReadMessages(client)
	go h.hub.WriteMessages(client)
}
