package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Timestamp string `json:"timestamp"`
}

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	wsConn *websocket.Conn
)

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	var err error
	wsConn, err = wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("could not upgrade: %s\n", err.Error())
		return
	}

	defer wsConn.Close()

	// event loop
	for {
		var msg Message

		err := wsConn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("error reading JSON: %s\n", err.Error())
			break
		}

		htmlContent, err := SendMessage(msg.Name, msg.Message)
		if err != nil {
			fmt.Printf("error generating HTML: %s\n", err.Error())
			break
		}

		err = wsConn.WriteMessage(websocket.TextMessage, []byte(htmlContent))
		if err != nil {
			fmt.Printf("error sending message: %s\n", err.Error())
			break
		}
	}
}

func SendMessage(name, message string) (string, error) {
	tmpl, err := template.ParseFiles("message.tmpl")
	if err != nil {
		return "", err
	}

	// Generate the timestamp in the desired format "d/m/yy hh:mm"
	timestamp := time.Now().Format("2/1/06 15:04")

	data := Message{
		Name:      name,
		Message:   message,
		Timestamp: timestamp,
	}

	var renderedHTML bytes.Buffer
	if err := tmpl.Execute(&renderedHTML, data); err != nil {
		return "", err
	}

	return renderedHTML.String(), nil
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ws", WsEndpoint)
	log.Fatal(http.ListenAndServe(":8000", router))
}
