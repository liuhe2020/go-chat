package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// client represents a single chatting user.

type client struct {
	// socket is the web socket for this client.
	socket *websocket.Conn
	// receive is a channel to receive messages from other clients.
	receive chan Message
	// room is the room this client is chatting in.
	room *room
}

type Message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg Message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.receive {
		timestamp := time.Now()
		htmlString := fmt.Sprintf(`<div id="chat_room" hx-swap-oob="beforeend"><p>%s - %s -%s<p></div>`, msg.Name, msg.Message, timestamp.Format("2/1/06 15:04"))
		htmlBytes := []byte(htmlString)
		err := c.socket.WriteMessage(websocket.TextMessage, htmlBytes)
		if err != nil {
			return
		}
	}
}
