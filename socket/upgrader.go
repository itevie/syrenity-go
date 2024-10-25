package socket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"syrenity/server/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Client) Read(db *sqlx.DB) {
	defer func() {
		c.server.unregister <- c
		c.Connection.Close()
	}()

	for {
		var msg Message
		err := c.Connection.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		switch msg.Type {
		case MTIdentify:
			var data struct {
				Token string `json:"token"`
			}

			err := json.Unmarshal([]byte(msg.Payload), &data)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			var user models.User
			authErr := models.UserFromAuth(db, data.Token, &user)

			if authErr != nil {
				fmt.Println(authErr.Error())
				return
			}

			c.User = &user

			helloData, _ := json.Marshal(user)
			c.Connection.WriteJSON(Message{
				Type:    MTHello,
				Payload: string(helloData),
			})
		}
	}
}

func (c *Client) Write(db *sqlx.DB) {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} else {
				err := c.Connection.WriteJSON(message)
				if err != nil {
					fmt.Println("Failed to write to WS: ", err)
					break
				}
			}
		case <-ticker.C:
			if err := c.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func NewClient(id string, conn *websocket.Conn, server *WebsocketServer) *Client {
	return &Client{ID: id, Connection: conn, send: make(chan Message, 256), server: server}
}

func Serve(c *gin.Context, s *WebsocketServer, db *sqlx.DB) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := NewClient(uuid.New().String(), ws, s)
	s.register <- client
	go client.Read(db)
	go client.Write(db)

	client.send <- Message{
		Type: MTAuthenticate,
	}
}
