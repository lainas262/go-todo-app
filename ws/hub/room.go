package hub

import (
	"context"
	"encoding/json"
	"example/todolist/ws/message"
	"log"
	"time"
)

type room struct {
	id               int64
	collaborators    map[*Client]bool
	registerClient   chan *Client
	unregisterClient chan *Client
	ctx              context.Context
	cancel           context.CancelFunc
	broadcast        chan message.Event
	hub              *Hub
}

func NewRoom(parentCtx context.Context, todoId int64, hub *Hub) *room {
	ctx, cancel := context.WithCancel(parentCtx)
	return &room{
		id:               todoId,
		collaborators:    make(map[*Client]bool),
		registerClient:   make(chan *Client),
		unregisterClient: make(chan *Client),
		ctx:              ctx,
		cancel:           cancel,
		broadcast:        make(chan message.Event),
		hub:              hub,
	}
}

func (r *room) Collaborators() map[message.Client]bool {
	out := make(map[message.Client]bool, len(r.collaborators))
	for id, val := range r.collaborators {
		out[id] = val
	}
	return out
}

func (r *room) RegisterClient(c *Client) {
	r.registerClient <- c
}

func (r *room) UnregisterClient(c *Client) {
	r.unregisterClient <- c
}
func (r *room) Broadcast(responseEvent message.Event) {
	r.broadcast <- responseEvent
}

func (r *room) Run() {
	logTickerForDebug := time.NewTicker(10 * time.Second)
	for {
		select {
		case client := <-r.registerClient:
			log.Println("Registered client in room!!!")
			r.collaborators[client] = true
		case client := <-r.unregisterClient:
			log.Println("Removed client from room!!!")
			delete(r.collaborators, client)
			if len(r.collaborators) == 0 {
				r.Close()
			}
		case message := <-r.broadcast:
			for client := range r.collaborators {
				if message.Client == client {
					continue
				}

				raw, err := json.Marshal(message)
				if err != nil {
					log.Println("Unable to marshal message")
					break
				}
				// log.Println(raw)
				select {
				case client.send <- &raw:
					//message successfully queued
				default:
					log.Printf("dropping message to slow client %v", client.userId)
				}
			}
			//notify clients
		case <-r.ctx.Done():
			return
		case <-logTickerForDebug.C:
			// log.Printf("Room Id: %v - Room Collaborators # %v", r.id, len(r.collaborators))
		}

	}
}

func (r *room) Close() {
	log.Println("Running Close room func")
	r.cancel()
	r.hub.unregisterRoom <- r
}
