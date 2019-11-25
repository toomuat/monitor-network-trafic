package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// type Message struct {
// 	Os      string `json:"os"`
// 	Counter uint64 `json:"counter"`
// }

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	randNum int

	// numString string
	// byteData  []byte
	// iphoneCounter  int = 0
	// androidCounter int = 0
	// windowsCounter int = 0

	// messagesMap map[string]*Message = map[string]*Message{
	// 	"Windows": &Message{"Windows", 0},
	// 	"Android": &Message{"Android", 0},
	// 	"iOS":     &Message{"iOS", 0},
	// }
	// messages []Message
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/send", exchangeData)

	log.Printf("Start HTTP server on localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func exchangeData(w http.ResponseWriter, r *http.Request) {
	upgrade, _ := upgrader.Upgrade(w, r, nil)

	rand.Seed(time.Now().UnixNano())

	for {
		messages = messages[:0] // clear slice
		for key, msg := range messagesMap {
			randNum = rand.Intn(100) // (0,100]
			messagesMap[key].Counter = uint64(randNum)
			messages = append(messages, *msg)
		}

		jsonBytes, err := json.Marshal(messages)
		if err != nil {
			log.Fatal(err)
		}
		jsonStr := string(jsonBytes)

		upgrade.WriteJSON(jsonStr)
		time.Sleep(1 * time.Second)
	}
}
