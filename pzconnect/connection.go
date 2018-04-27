package pzconnect

import (
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/goinggo/tracelog"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 120 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageBytes = 1024
)

// IConnection is an interface to connection
type IConnection interface {
	Start()
	SendMessage(body []byte)
	SendClose(closeCode int, text string)
	GetRequest() *http.Request
	SetShouldRemoveFromHub(value bool)
}

type message struct {
	messageType int
	body        []byte
}

type client struct {
	ClientID    uint64
	QueryParams url.Values
}

type connection struct {
	ws     *websocket.Conn
	req    *http.Request
	client *client
	send   chan message
	quit   chan struct{}

	mu                  sync.Mutex
	shouldRemoveFromHub bool
}

func newConnection(ws *websocket.Conn, r *http.Request, params url.Values) *connection {
	if ws == nil {
		tracelog.Trace("pzconnect", "NewConnection", "nil value passed as ws to NewConnection")
		return nil
	}

	client := &client{
		ClientID:    0,
		QueryParams: params,
	}
	connection := &connection{
		ws:                  ws,
		req:                 r,
		client:              client,
		send:                make(chan message, 256),
		quit:                make(chan struct{}),
		shouldRemoveFromHub: true,
	}

	return connection
}

func (c *connection) SetShouldRemoveFromHub(value bool) {
	c.mu.Lock()
	c.shouldRemoveFromHub = value
	c.mu.Unlock()
}

func (c *connection) SendMessage(body []byte) {
	c.send <- message{
		websocket.BinaryMessage,
		body,
	}
}

func (c *connection) SendClose(closeCode int, text string) {
	c.send <- message{
		websocket.CloseMessage,
		websocket.FormatCloseMessage(closeCode, text),
	}
}

func (c *connection) GetRequest() *http.Request {
	return c.req
}

func (c *connection) Start() {
	go c.writer()
	go c.reader()
}

func (c *connection) writer() {
	defer func() {
		tracelog.Info("pzconnect", "writer", "Writer stopped %#x", c.client.ClientID)
	}()

	for {
		select {
		case <-c.quit:
			return

		case message := <-c.send:
			if err := c.write(message.messageType, message.body); err != nil {
				return
			}
			if message.messageType == websocket.CloseMessage {
				return
			}
		}
	}
}

func (c *connection) write(messageType int, payload []byte) error {
	err := c.ws.WriteMessage(messageType, payload)
	if err != nil {
		tracelog.Errorf(err, "pzconnect", "write", "Error sending message:")
	}
	return err
}

func (c *connection) reader() {
	defer tracelog.Info("pzconnect", "reader", "Reader stopped %#x", c.client.ClientID)
	defer c.ws.Close()
	defer close(c.quit)
	defer func() {
		tracelog.Info("pzconnect", "reader", "client id %#x, shouldRemoveFromHub %d", c.client.ClientID, c.shouldRemoveFromHub)
		if c.shouldRemoveFromHub == true {
			tracelog.Info("pzconnect", "reader", "deleting from hub %#x", c.client.ClientID)
			clientsHub.Delete(c.client.ClientID)
		} else {
			tracelog.Info("pzconnect", "reader", "not deleting from hub, due to connection is an orphan %#x", c.client.ClientID)
		}
	}()

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			switch err.(type) {
			case *websocket.CloseError:
				closeErr := err.(*websocket.CloseError)
				tracelog.Errorf(closeErr, "pzconnect", "reader", "Socket closed %#x", c.client.ClientID)
			default:
				tracelog.Errorf(err, "pzconnect", "reader", "Error in reader %#x", c.client.ClientID)
			}
			return
		}

		go c.handleMessage(message)
	}
}

func (c *connection) handleMessage(message []byte) {
	tracelog.Info("pzconnect", "handleMessage", "got message from %#x", c.client.ClientID)
	handleRequest(c, message)
}
