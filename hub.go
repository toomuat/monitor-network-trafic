package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan string
	register   chan *Client
	unregister chan *Client
}

var (
	clientNum uint64

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan string),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Printf("new client [%d] joined\n", client.id)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				fmt.Printf("client [%d] left\n", client.id)
				clientNum--
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
					// fmt.Printf("%d: %v\n", client.id, message)
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			iosCounter = 0
			androidCounter = 0
			windowsCounter = 0
		}
	}
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrade, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	clientNum++
	client := &Client{
		hub:  h,
		conn: upgrade,
		send: make(chan string),
		id:   clientNum,
	}
	client.hub.register <- client

	go client.broardcastCounter()
	go client.sendCounter()
	go client.detectDisconnection()
}
