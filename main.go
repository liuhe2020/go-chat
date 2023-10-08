package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Message struct {
	Name      string `json:"name"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type ChatHistory struct {
	Messages []Message `json:"messages"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn *websocket.Conn
)

var chatHistory ChatHistory

func init() {
	// Load existing messages from data.json
	loadChatHistory()
}

func loadChatHistory() {
	file, err := os.ReadFile("data.json")
	if err != nil {
		log.Println("Error reading data.json:", err)
		return
	}

	err = json.Unmarshal(file, &chatHistory)
	if err != nil {
		log.Println("Error unmarshalling data.json:", err)
		return
	}
}

func saveChatHistory() {
	data, err := json.MarshalIndent(chatHistory, "", "  ")
	if err != nil {
		log.Println("Error marshalling chat history:", err)
		return
	}

	err = os.WriteFile("data.json", data, 0644)
	if err != nil {
		log.Println("Error writing data.json:", err)
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	loadChatHistory()
	tmpl.Execute(w, chatHistory)
}

var clients = make(map[*websocket.Conn]bool) // Connected clients

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("could not upgrade: %s\n", err.Error())
		return
	}
	defer conn.Close()

	// Add the new client to the list of connected clients
	clients[conn] = true

	// Event loop for reading messages from the WebSocket
	for {
		var msg Message

		// Read a message from the WebSocket
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("WebSocket connection closed by client")
			} else {
				log.Printf("error reading JSON: %s\n", err.Error())
			}
			break
		}

		timestamp := time.Now()

		// Construct the HTML message using the data from the Message object
		htmlString := fmt.Sprintf(`<div id="chat_room" hx-swap-oob="beforeend"><p>%s - %s -%s<p></div>`, msg.Name, msg.Message, timestamp.Format("2/1/06 15:04"))
		htmlBytes := []byte(htmlString)

		// Broadcast HTML content to all clients
		for client := range clients {
			err = client.WriteMessage(websocket.TextMessage, htmlBytes)
			if err != nil {
				log.Printf("error writing to WebSocket: %s\n", err.Error())
				// Remove the disconnected client from the list
				delete(clients, client)
				break
			}
		}

		// Convert the timestamp to ISO 8601 string for storage in data.json
		msg.Timestamp = timestamp.UTC().Format(time.RFC3339)

		chatHistory.Messages = append(chatHistory.Messages, msg)

		// Save the updated chat history to the file
		saveChatHistory()
	}

	// Remove the client from the list when the loop ends
	delete(clients, conn)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", handler)
	router.HandleFunc("/ws", wsHandler)

	log.Fatal(http.ListenAndServe(":8000", router))
}
