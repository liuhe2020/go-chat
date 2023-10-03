package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

type Message struct {
	Name      string  `json:"name"`
	Message   string  `json:"message"`
	Timestamp string  `json:"timestamp"`
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New connection:", ws.RemoteAddr())

	connectionMessage := "⚡ Connected to Go-Chat server"
	if err := websocket.Message.Send(ws, connectionMessage); err != nil {
		fmt.Println("Error sending welcome message:", err)
		return
	}

	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	// Load existing data from data.json, if any
	// existingData, err := loadData()
	// if err != nil {
	// 	fmt.Println("Error loading existing data:", err)
	// 	return
	// }

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

		fmt.Printf("Received message from %s: %s\n", msg.Name, msg.Message)

		// Get the current timestamp
		msg.Timestamp = time.Now().Format(time.RFC3339)

		// Append the new message to existing data
		// existingData = append(existingData, msg)

		// Write the updated data to data.json
		// err = saveData(existingData)
		// if err != nil {
		// 	fmt.Println("Error saving data:", err)
		// 	return
		// }

		fmt.Println("Message saved to data.json")

		// Render the message template
		renderedMessage, err := renderMessageTemplate(&msg)
		if err != nil {
			fmt.Println("Error rendering message template:", err)
			continue
		}

		// Broadcast the HTML content
		s.broadcast([]byte(renderedMessage))
	}
}

func loadData() ([]Message, error) {
	var data []Message

	// Read data from data.json
	file, err := os.Open("data.json")
	if err != nil {
		return data, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func saveData(data []Message) error {
	// Write data to data.json
	file, err := os.Create("data.json")
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(data)
	if err != nil {
		return err
	}

	return nil
}


func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("Write error:", err)
			}
		}(ws)
	}
}

func renderMessageTemplate(msg *Message) (string, error) {
	tmpl := `
		<p>{{.Name}}</p>
		<p>{{.Message}}</p>
		<p>{{.Timestamp}}</p>`

	messageTemplate, err := template.New("messageTemplate").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var renderedMessage bytes.Buffer
	err = messageTemplate.Execute(&renderedMessage, msg)
	if err != nil {
		return "", err
	}

	return renderedMessage.String(), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Yo Momma")
}

func main() {
	server := NewServer()
	mux := http.NewServeMux()
	mux.Handle("/ws", websocket.Handler(server.handleWS))
	mux.HandleFunc("/", handler)

	fmt.Println("Server is running on :8000")
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
