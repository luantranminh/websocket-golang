package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

//Client contain unique ID and hold their connection
type Client struct {
	ID         string
	Connection *websocket.Conn
}

// PubSub relation between clients and their subcription
type PubSub struct {
	Subscription []Subscription
}

//Message define structure of message send through websocket
type Message struct {
	GroupID string `json:"group_id"`
	Lat     int    `json:"lat"`
	Lng     int    `json:"lng"`
	UserID  string `json:"user_id"`
}
type Subscription struct {
	Clients []Client
	GroupID string
}

func (ps *PubSub) getSubcriptions(groupID string) []Client {

	for _, el := range ps.Subscription {
		if el.GroupID == groupID {
			return el.Clients
		}
	}
	return []Client{}
}

func (ps *PubSub) isSubscribed(client Client, subs Subscription) bool {
	for _, el := range subs.Clients {
		if el.ID == client.ID {
			return true
		}
	}

	return false
}

func (ps *PubSub) subscribe(groupID string, client Client) {
	for index, el := range ps.Subscription {
		if el.GroupID == groupID {
			if ps.isSubscribed(client, ps.Subscription[index]) {
				return
			}
			ps.Subscription[index].Clients = append(ps.Subscription[index].Clients, client)
			return
		}
	}

	newSubscription := Subscription{
		Clients: []Client{
			client,
		},
		GroupID: groupID,
	}
	ps.Subscription = append(ps.Subscription, newSubscription)
}

func (ps *PubSub) HandleReceiveMessage(payload []byte, client Client) {

	msg := Message{}

	err := json.Unmarshal(payload, &msg)
	if err != nil {
		fmt.Println("This is not correct message payload")
	}

	client.ID = msg.UserID
	ps.subscribe(msg.GroupID, client)

	clientsHasSubscribed := ps.getSubcriptions(msg.GroupID)

	for _, client := range clientsHasSubscribed {
		if msg.UserID != client.ID {
			client.Connection.WriteMessage(websocket.TextMessage, payload)
		}
	}

}
