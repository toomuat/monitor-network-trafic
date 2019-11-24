package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Os      string `json:"os"`
	Counter uint64 `json:"counter"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var randNum int
var numString string
var byteData []byte
var iphoneCounter int = 0
var androidCounter int = 0
var windowsCounter int = 0
var messagesMap map[string]*Message = map[string]*Message{
	"Windows": &Message{"Windows", 0},
	"Android": &Message{"Android", 0},
	"iOS":     &Message{"iOS", 0},
}
var messages []Message

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
			// fmt.Printf("%v: {%v, %v}\n", key, (*msg).Counter, (*msg).Os)
			// fmt.Printf("%v: {%v, %v}\n", key, msg.Counter, msg.Os) // ok
			messages = append(messages, *msg)
		}

		jsonBytes, err := json.Marshal(messages)
		if err != nil {
			log.Fatal(err)
		}
		jsonStr := string(jsonBytes)

		// numString = strconv.FormatUint(uint64(randNum), 10)
		// fmt.Printf("%v\n", numString)
		// byteData = []byte(numString)
		// upgrade.WriteMessage(websocket.TextMessage, byteData)
		upgrade.WriteJSON(jsonStr)
		time.Sleep(1 * time.Second)
	}
}
