package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var randNum int

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/send", exchangeData)

	// only work on linux
	// go func(){
	// 	capturePacket()
	// }()

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
