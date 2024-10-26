// Mostly stolen / inspred from https://github.com/tinkerbaj/chat-websocket-gin
package socket

import (
	"fmt"
	"strings"
	"syrenity/server/models"

	"github.com/gorilla/websocket"
)

const (
	MTAuthenticate = "Authenticate"
	MTIdentify     = "Identify"
	MTHello        = "Hello"
	MTClose        = "Close"
	MTError        = "Error"

	MTDMessageCreate = "DispatchMessageCreate"
)

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type Client struct {
	ID         string
	Connection *websocket.Conn
	User       *models.User
	send       chan Message
	server     *WebsocketServer
}

type WebsocketServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
}

func NewServer() *WebsocketServer {
	return &WebsocketServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

func (s *WebsocketServer) RegisterClient(client *Client) {
	s.clients[client] = true

	fmt.Println("Size of clients: ", len(s.clients))
}

func (s *WebsocketServer) RemoveClient(client *Client) {
	delete(s.clients, client)
	close(client.send)
	fmt.Println("Removed client")
}

func (s *WebsocketServer) Run() {
	for {
		select {
		case client := <-s.register:
			s.RegisterClient(client)
		case client := <-s.unregister:
			s.RemoveClient(client)
		case message := <-s.broadcast:
			s.HandleMessage(message)
		}
	}
}

func (s *WebsocketServer) HandleMessage(message Message) {
	for client := range s.clients {
		if !strings.HasPrefix(message.Type, "Dispatch") && client.User == nil {
			continue
		}

		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(s.clients, client)
		}
	}
}
