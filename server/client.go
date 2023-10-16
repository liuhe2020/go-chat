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
	// keep track of user who sent the message
	username string
}

type Message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (c *client) read() {
	defer c.socket.Close()
	// Continue reading messages
	for {
		var msg Message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		// set or update the username
		c.username = msg.Name
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.receive {
		timestamp := time.Now()
		var cssClasses string

		if msg.Name == c.username {
			cssClasses = "bg-indigo-500 text-white self-end"
		} else {
			cssClasses = "bg-[#f0f0f1] text-gray-900"
		}

		htmlString := fmt.Sprintf(`<div
          id="chat_room" hx-swap-oob="beforeend"
          class="py-4 px-2 flex flex-1 flex-col gap-y-4 overflow-y-auto scrollbar-thumb-blue scrollbar-thumb-rounded scrollbar-track-blue-lighter scrollbar-w-2 scrolling-touch"
        >
          <div class="flex flex-col max-w-xs px-4 py-2.5 rounded-md inline-block %s">
            <div class="flex gap-x-2 pt-0.5 pb-1">
              <span class="align-bottom font-bold leading-3">%s</span>
              <span class="align-bottom text-xs font-medium">%s</span>
            </div>
            <span class="leading-5">%s</span>
          </div></div>`, cssClasses, msg.Name, timestamp.Format("2/1/06 15:04"), msg.Message)

		htmlBytes := []byte(htmlString)
		err := c.socket.WriteMessage(websocket.TextMessage, htmlBytes)
		if err != nil {
			return
		}
	}
}
