package service

import (
	"context"
	"log"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type ClientMessage struct {
	UserId  int64
	Message any
}

type WebsocketService struct {
	Connections map[*Client]bool
	Register    chan *Client
	Unregister  chan *Client
	Incoming    chan *ClientMessage
	Outgoing    chan *ClientMessage
}

type Client struct {
	Conn   *websocket.Conn
	UserId int64
	Send   chan *ClientMessage
}

func CreateWebsocketService() *WebsocketService {
	return &WebsocketService{
		Connections: make(map[*Client]bool),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Incoming:    make(chan *ClientMessage),
		Outgoing:    make(chan *ClientMessage),
	}
}

func (s *WebsocketService) StartHub() {
	for {
		select {
		case client := <-s.Register:
			log.Println("Register user: ", client)
			s.Connections[client] = true
			log.Println("Registered user: ", client)
		case client := <-s.Unregister:
			if _, ok := s.Connections[client]; ok {
				log.Println("Unregister user: ", client)
				delete(s.Connections, client)
				log.Println("Unregistered user: ", client)
			}
		case message := <-s.Incoming:
			log.Printf("Broadcast user message - userId: %v, message: %v", message.UserId, message.Message)
			for client := range s.Connections {
				client.Send <- message
			}
		}

	}
}

func (s *WebsocketService) ReadMessages(client *Client) {
	defer func() {
		client.Conn.CloseNow()
		s.Unregister <- client
	}()

	for {

		var v any
		err := wsjson.Read(context.Background(), client.Conn, &v)
		if err != nil {
			log.Printf("Connection Gone %v", err)
			break
		}
		s.Incoming <- &ClientMessage{
			UserId:  client.UserId,
			Message: v,
		}
	}
}

func (s *WebsocketService) WriteMessages(client *Client) {
	defer func() {
		client.Conn.CloseNow()
		log.Println("Calling defer")
	}()

	for {
		select {
		case message, ok := <-client.Send:

			if !ok {
				return
			}
			if client.UserId != message.UserId {
				log.Printf("Test write message - userId: %v, message: %v", message.UserId, message.Message)
				err := wsjson.Write(context.Background(), client.Conn, message.Message)
				if err != nil {
					log.Printf("Unable to write message: %v", err)
					return
				}
			}
		default:
			client.Conn.Ping(context.Background())
		}

	}
}
