package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//Message define structure of message send through websocket
type Message struct {
	// Name    string `json:"name"`
	// Title   string `json:"title"`
	Content string `json:"content"`
}

var upgrader = websocket.Upgrader{}
var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

func main() {
	//Create a simple file

	fs := http.FileServer(http.Dir(""))
	http.Handle("/", fs)

	//configure websocket route
	http.HandleFunc("/ws", handleConnections)

	//start listening for incoming chat message
	go handleMessages()

	log.Println("http server started on :4321")
	http.ListenAndServe(":4321", nil)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	//Register new client

	clients[ws] = true

	for {
		var msg Message
		//Read new message as JSON and map it to Message object

		err := ws.ReadJSON(&msg)
		fmt.Println(msg)
		if err != nil {
			log.Printf("error when handle connections: %v ", err)
			delete(clients, ws)
			break
		}

		// Send newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg.Content)

			if err != nil {
				log.Printf("error when handle messages: %v ", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
