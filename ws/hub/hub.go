package hub

import (
	"context"
	"encoding/json"
	"example/todolist/api"
	"time"

	"example/todolist/ws/handler"
	"example/todolist/ws/message"
	"log"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/redis/go-redis/v9"
)

type Hub struct {
	connections      map[*Client]bool
	registerClient   chan *Client
	unregisterClient chan *Client
	rooms            map[int64]*room
	registerRoom     chan *registerRequest
	unregisterRoom   chan *room
	incoming         chan *ClientMessage
	// Outgoing    chan *Event
	redis  *redis.Client
	router *Router
}

func NewHub(redis *redis.Client) *Hub {
	return &Hub{
		connections:      make(map[*Client]bool),
		registerClient:   make(chan *Client),
		unregisterClient: make(chan *Client),
		rooms:            make(map[int64]*room),
		registerRoom:     make(chan *registerRequest),
		unregisterRoom:   make(chan *room),
		incoming:         make(chan *ClientMessage),
		// Outgoing:    make(chan *Event),
		redis:  redis,
		router: NewEventRouter(),
	}
}

func (h *Hub) RegisterClient(c *Client) {
	h.registerClient <- c
}

func (h *Hub) UnregisterClient(c *Client) {
	h.unregisterClient <- c
}

func (h *Hub) RegisterRoom(r *registerRequest) {
	h.registerRoom <- r
}

func (h *Hub) UnregisterRoom(r *room) {
	h.unregisterRoom <- r
}

func (h *Hub) Incoming(m *ClientMessage) {
	h.incoming <- m
}

func (h *Hub) GetRoom(id int64) message.Room {
	var room message.Room = h.rooms[id]
	return room

}

func (h *Hub) RegisterRoutes() {
	h.router.Register("open_todo", &handler.OpenTodoHandler{Hub: h})
	h.router.Register("update_todo", &handler.UpdateTodoHandler{Hub: h})
}

func (h *Hub) Run() {
	logTickerForDebug := time.NewTicker(10 * time.Second)
	for {
		select {
		case client := <-h.registerClient:
			// log.Println("Register user: ", client)
			h.connections[client] = true
			// log.Println("Registered user: ", client)
		case client := <-h.unregisterClient:
			// log.Println("Unregister user: ", client)
			delete(h.connections, client)
			for _, room := range h.rooms {
				room.unregisterClient <- client
			}
			// log.Println("Unregistered user: ", client)
		case req := <-h.registerRoom:
			// create room based on todo item
			room, ok := h.rooms[req.todo.Id]
			if !ok {
				room = NewRoom(context.Background(), req.todo.Id, h)
				go room.Run()
				h.rooms[room.id] = room
			}
			// register client here after creating room
			room.registerClient <- req.client
		case room := <-h.unregisterRoom:
			log.Println("call unregister rooms need to check if there's users in it before")
			delete(h.rooms, room.id)
		case message := <-h.incoming:
			//could leverage workers perhaps at some point?
			go h.processIncomingMessage(*message)
		case <-logTickerForDebug.C:
			// log.Printf("Number of rooms: %v", len(h.rooms))

		}

	}
}

func (h *Hub) ReadMessages(client *Client) {
	defer func() {
		client.cancel()
		client.conn.CloseNow()
		h.UnregisterClient(client)
	}()

	for {
		_, b, err := client.conn.Read(client.ctx)
		if err != nil {
			log.Printf("Connection Gone %v", err)
			return
		}
		h.Incoming(&ClientMessage{
			Client: client,
			Raw:    b,
		})
	}
}

func (h *Hub) WriteMessages(client *Client) {
	pingTicker := time.NewTicker(30 * time.Second)
	defer func() {
		client.conn.CloseNow()
		pingTicker.Stop()
	}()

	for {
		select {
		case message, ok := <-client.send:

			if !ok {
				return
			}

			err := client.conn.Write(client.ctx, websocket.MessageText, *message)
			if err != nil {
				log.Printf("Unable to write message - reason: %v", err)
				return
			}
		case <-pingTicker.C:
			if err := client.conn.Ping(context.Background()); err != nil {
				log.Printf("Ping Error: %v", err)
			}
		case <-client.ctx.Done():
			log.Println("Client Context cancelled")
			return
		}
	}
}

func (h *Hub) processIncomingMessage(clientMessage ClientMessage) {
	var e message.Event
	err := json.Unmarshal(clientMessage.Raw, &e)
	if err != nil {
		log.Println("failed to unmarshal json: ", err)
		return
	}
	if e.Name == "" {
		log.Println("'event' can not be empty or undefined")
		return
	}
	if handler, ok := h.router.rules[e.Name]; ok {
		e.Client = clientMessage.Client

		handler.Handle(e.Client.Ctx(), e)
	}
}

type registerRequest struct {
	client *Client
	todo   *api.ResponseTodoWithCollaborators
}

func (h *Hub) CreateClientRooms(c *Client) {
	//perhaps see if we can create new context to pass down to http calls to cancel them if client disconnects
	var wg sync.WaitGroup
	// go routine to fetch client todos, and start
	// creating rooms for each todo if they don't exist
	client := &http.Client{}

	var todos []api.ResponseTodoWithCollaborators
	var collaborationTodos []api.ResponseTodoWithCollaborators

	wg.Go(func() {
		req, err := http.NewRequestWithContext(c.ctx, "GET", "http://localhost:8080/api/todos", nil)

		if err != nil {
			log.Println(err)
			return
		}

		req.Header.Set("authorization", "Bearer "+c.token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&todos)
		if err != nil {
			log.Println("error decoding json: ", err)
		}
	})

	wg.Go(func() {
		req, err := http.NewRequestWithContext(c.ctx, "GET", "http://localhost:8080/api/collaborations", nil)

		if err != nil {
			log.Println(err)
			return
		}

		req.Header.Set("authorization", "Bearer "+c.token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&collaborationTodos)
		if err != nil {
			log.Println("error decoding json: ", err)
		}
	})

	wg.Wait()

	items := append(todos, collaborationTodos...)

	defer log.Println("Goroutine exited to cleanup?? :P")

	for _, item := range items {
		select {
		case <-c.ctx.Done():
			log.Println("canceled")
			return
		default:
			log.Println("Sending register Request")
			h.registerRoom <- &registerRequest{
				client: c,
				todo:   &item,
			}
		}

	}
}
