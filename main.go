package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
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

func (s *Server) readLoop(ws *websocket.Conn){
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error:", err)
			continue
		}
		msg := buf[:n]
		s.broadcast(msg)
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