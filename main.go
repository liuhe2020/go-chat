package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn){
	fmt.Println("new connection:", ws.RemoteAddr())

	connectionMessage := "⚡Connected to Go-Chat server"
	if err := websocket.Message.Send(ws, connectionMessage); err != nil {
		fmt.Println("error sending welcome message:", err)
		return
	}

	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	for {
		var msg Message
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Read error:", err)
			continue
		}

		fmt.Printf("Received message from %s: %s\n", msg.Username, msg.Message)

		// Marshal the message into JSON bytes
		jsonData, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error marshalling message:", err)
			continue
		}

		// Write the JSON data to a file
	err = os.WriteFile("data.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Data written to data.json")

		// Broadcast the JSON data
		s.broadcast(jsonData)
	}
}


func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error:", err)
			}
		}(ws)
	}
}

func main() {
	server := NewServer()
	fmt.Println("Websocket server is running on :1337")
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.ListenAndServe(":1337", nil)
}