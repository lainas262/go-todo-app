package hub

import (
	"context"

	"github.com/coder/websocket"
)

type Client struct {
	token  string
	conn   *websocket.Conn
	userId int64
	send   chan *[]byte
	ctx    context.Context
	cancel context.CancelFunc
}

type ClientMessage struct {
	Client *Client
	Raw    []byte
	Room   *room
}

func NewClient(token string, userId int64, wsConn *websocket.Conn) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		token:  token,
		conn:   wsConn,
		userId: userId,
		send:   make(chan *[]byte, 256),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (c *Client) Send(msg []byte) {
	// we'll see about the actual implementation
}

func (c *Client) Ctx() context.Context {
	return c.ctx
}

func (c *Client) Conn() *websocket.Conn {
	return c.conn
}
