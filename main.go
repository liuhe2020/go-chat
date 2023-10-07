package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Message struct {
	Name      string `json:"name"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	wsConn *websocket.Conn
)

func Handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))

	messages := map[string][]Message{
		"Messages": {
			{Name: "The Godfather", Message: "Francis Ford Coppola", Timestamp: "10/10/10"},
			{Name: "Blade Runner", Message: "Ridley Scott", Timestamp: "10/10/10"},
			{Name: "The Thing", Message: "John Carpenter", Timestamp: "10/10/10"},
		},
	}

	tmpl.Execute(w, messages)

}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// Upgrade the HTTP connection to a WebSocket connection
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("could not upgrade: %s\n", err.Error())
		return
	}
	defer wsConn.Close()

	// Event loop for reading messages from the WebSocket
	for {
		var msg Message

		// Read a message from the WebSocket
		err := wsConn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("WebSocket connection closed by client")
			} else {
				log.Printf("error reading JSON: %s\n", err.Error())
			}
			break
		}

		timestamp := time.Now().Format("2/1/06 15:04")

		// Construct the HTML message using the data from the Message object
		htmlString := fmt.Sprintf(`<div id="chat_room" hx-swap-oob="beforeend"><p>%s - %s -%s<p></div>`, msg.Name, msg.Message, timestamp)
		htmlBytes := []byte(htmlString)

		// Send HTML content to the client over the WebSocket
		err = wsConn.WriteMessage(websocket.TextMessage, htmlBytes)
		if err != nil {
			log.Printf("error writing to WebSocket: %s\n", err.Error())
			break
		}
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", Handler)
	router.HandleFunc("/ws", WsEndpoint)

	log.Fatal(http.ListenAndServe(":8000", router))
}
