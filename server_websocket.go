package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var randNum int
var numString string
var byteData []byte

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	// http.HandleFunc("/chart", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./static/cht.html")
	// })
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
		randNum = rand.Intn(100)
		numString = strconv.FormatUint(uint64(randNum), 10)
		fmt.Printf("%v\n", numString)
		byteData = []byte(numString)
		upgrade.WriteMessage(websocket.TextMessage, byteData)
		time.Sleep(1 * time.Second)
	}
}
