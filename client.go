package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Message : packet counter of each OS
type Message struct {
	Os      string `json:"os"`
	Counter uint64 `json:"counter"`
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan string
	id   uint64
}

var(
	messages = []Message{
		{Os: "Windows"},
		{Os: "Android"},
		{Os: "iOS"},
	}
)

func (c *Client) broardcastCounter() {
	for {
		// send number of packet every 1 second
		time.Sleep(1 * time.Second)

		messages[0].Counter = windowsCounter
		messages[1].Counter = androidCounter
		messages[2].Counter = iosCounter

		jsonBytes, err := json.Marshal(messages)
		if err != nil {
			log.Fatal(err)
		}

		c.hub.broadcast <- string(jsonBytes)
	}
}

func (c *Client) sendCounter() {
	for {
		select {
		case message, _ := <-c.send:
			err := c.conn.WriteJSON(message)
			if err != nil {
				c.hub.unregister <- c
				c.conn.Close()
				return
				// break
			}
		}
	}
	// c.conn.Close()
}

func (c *Client) detectDisconnection() {
	_, _, err := c.conn.ReadMessage()
	if err != nil {
		c.hub.unregister <- c
		c.conn.Close()
	}
}
