package handler

import (
	"example/todolist/middleware"
	"example/todolist/response"
	service "example/todolist/ws/service"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type WebsocketHandler struct {
	service *service.WebsocketService
}

func CreateWebsocketHandler(service *service.WebsocketService) *WebsocketHandler {
	return &WebsocketHandler{service: service}
}

func (h *WebsocketHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIdKey).(int64)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	log.Printf("Websocket Subscribe - Request received: %v", userId)

	wsConn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"http://localhost:3000"},
	})

	if err != nil {
		log.Printf("Websocket accept err: %v", err)
		return
	}

	clientConn := &service.Client{
		UserId: userId,
		Conn:   wsConn,
		Send:   make(chan *service.ClientMessage),
	}

	h.service.Register <- clientConn
	go h.service.ReadMessages(clientConn)
	go h.service.WriteMessages(clientConn)
}
