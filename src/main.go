package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var upgrader = websocket.Upgrader{}

// var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message) // broadcast channel
var clients = PubSub{}

func autoID() string {
	return uuid.Must(uuid.NewV4()).String()
}

func main() {
	//Create a simple file

	fs := http.FileServer(http.Dir(""))
	http.Handle("/", fs)

	//configure websocket route
	http.HandleFunc("/ws", handleConnections)

	//start listening for incoming chat message
	// go handleMessages()

	log.Println("http server started on :4321")
	http.ListenAndServe(":4321", nil)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()

	//Register new client

	newClient := Client{
		Connection: conn,
	}

	// newClient := Client{
	// 	ID:         "1",
	// 	Connection: ws,
	// }
	// clients1.AddClient(newClient)

	for {
		// var msg Message
		//Read new message as JSON and map it to Message object
		_, p, err := conn.ReadMessage()

		// err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error when handle connections: %v ", err)
			break
		}
		clients.HandleReceiveMessage(p, newClient)
		// fmt.Println(p)
		// Send newly received message to the broadcast channel
		// broadcast <- msg
	}
}

// func handleMessages() {
// 	for {
// 		// Grab the next message from the broadcast channel
// 		msg := <-broadcast

// 		for client := range clients {
// 			err := client.WriteJSON(msg)

// 			if err != nil {
// 				log.Printf("error when handle messages: %v ", err)
// 				client.Close()
// 				delete(clients, client)
// 			}
// 		}
// 	}
// }
