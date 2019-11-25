package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Message struct {
	Os      string `json:"os"`
	Counter uint64 `json:"counter"`
}

var (
	windowsCounter uint64 = 0
	androidCounter uint64 = 0
	iosCounter     uint64 = 0

	messagesMap map[string]*Message = map[string]*Message{
		"Windows": &Message{"Windows", 0},
		"Android": &Message{"Android", 0},
		"iOS":     &Message{"iOS", 0},
	}
	messages  []Message
	messages2 []*Message = []*Message{
		&Message{"Windows", 0},
		&Message{"Android", 0},
		&Message{"iOS", 0},
	}
)

func sendCounter(w http.ResponseWriter, r *http.Request) {
	upgrade, _ := upgrader.Upgrade(w, r, nil)

	go packetCounter()
	for {
		// send number of packet every 1 second
		time.Sleep(1 * time.Second)

		messages = messages[:0] // clear slice
		messages[0].Counter = windowsCounter
		messages[1].Counter = androidCounter
		messages[2].Counter = iosCounter

		jsonBytes, err := json.Marshal(messages)
		if err != nil {
			log.Fatal(err)
		}
		jsonStr := string(jsonBytes)

		upgrade.WriteJSON(jsonStr)

		iosCounter = 0
		androidCounter = 0
		windowsCounter = 0
	}
}

func packetCounter() {
	for {
		// count packet of each OS
		time.Sleep(1 * time.Microsecond)
	}
}
