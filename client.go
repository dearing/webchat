package main

import (
	"code.google.com/p/go.net/websocket"
	. "fmt"
	"log"
)

type Client struct {
	ID       string
	Nickname string
	Socket   *websocket.Conn
}

type ClientJSON struct {
	Action string
	Data   string
	Origin string
}

func (client *Client) mesg(action, data string) *Message {
	return &Message{Client: client, JSON: &ClientJSON{Action: action, Data: data, Origin: client.Nickname}}
}

func WebsocketHandler(ws *websocket.Conn) {

	client := &Client{ID: <-Router.GetID, Socket: ws}

	// Getting a nickname off the bat is our basic handshake here...
	var q ClientJSON
	if err := websocket.JSON.Receive(client.Socket, &q); err != nil {
		return
	}
	if q.Action != "nickname" {
		if config.Verbose {
			log.Printf("[%s] failed handshake.\n", client.ID, client.Nickname)
		}
		return
	}
	client.Nickname = q.Data

	// add out client
	Router.Add <- client

	// Defer cleanup
	defer func() {
		Router.Remove <- client
		Router.Echo <- client.mesg("inform", Sprintf("%s left.", client.Nickname))
		if config.Verbose {
			log.Printf("[%s]:%s disconnected.\n", client.ID, client.Nickname)
		}
	}()

	if config.Verbose {
		log.Printf("[%s]:%s connected.\n", client.ID, client.Nickname)
	}
	Router.Echo <- client.mesg("inform", Sprintf("%s joined.", client.Nickname))

	// From here out we loop building JSON data and acting on it
	for {

		var q ClientJSON

		if err := websocket.JSON.Receive(client.Socket, &q); err != nil {
			break
		}

		if config.Verbose {
			println(Sprintf("%s:%s => %s => %s", client.ID, client.Nickname, q.Action, q.Data))
		}

		switch q.Action {

		case "disconnect":
			return

		case "nickname":
			old_nick := client.Nickname
			client.Nickname = q.Data
			Router.Echo <- client.mesg("inform", Sprintf("%s is now %s.", old_nick, client.Nickname))

		case "message":
			Router.Echo <- client.mesg("message", q.Data)
		}
	}
}
