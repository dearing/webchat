package main

import (
	"code.google.com/p/go.net/websocket"
	"crypto/sha1"
	"fmt"
	"time"
)

type router struct {
	Add       chan *Client
	Remove    chan *Client
	Broadcast chan *Message
	Echo      chan *Message
	Send      chan *Message
	Clients   map[*Client]string
	GetID     chan string
}

var Router = router{
	Add:       make(chan *Client),
	Remove:    make(chan *Client),
	Broadcast: make(chan *Message),
	Echo:      make(chan *Message),
	Send:      make(chan *Message),
	Clients:   make(map[*Client]string),
	GetID:     make(chan string),
}

type Message struct {
	Client *Client
	JSON   *ClientJSON
}

func (Router *router) HandleClients() {

	// provides a new uid on request
	go func() {
		h := sha1.New()
		for {
			h.Write([]byte(time.Now().String()))
			Router.GetID <- fmt.Sprintf("%X", h.Sum(nil))
		}
	}()

	for {
		select {

		// ADD
		case client := <-Router.Add:
			Router.Clients[client] = client.ID

		// REMOVE
		case client := <-Router.Remove:
			client.Socket.Close()
			delete(Router.Clients, client)

		// BROADCAST
		case message := <-Router.Echo:
			for peer, _ := range Router.Clients {
				websocket.JSON.Send(peer.Socket, message.JSON)
			}

		// BROADCAST (ALL !CLIENT)
		case message := <-Router.Broadcast:
			for peer, key := range Router.Clients {
				if peer != message.Client {
					if key == message.Client.ID {
						websocket.JSON.Send(peer.Socket, message.JSON)
					}
				}
			}

		// SEND
		case message := <-Router.Send:
			websocket.JSON.Send(message.Client.Socket, message.JSON)
		}
	}
}
