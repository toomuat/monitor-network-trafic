package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <interface>\n", filepath.Base(os.Args[0]))
		return
	}

	// create file to save log of the computer commected to the network
	fd, err := os.Create("srcMac-OS.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fd.Close()
	log.Printf("Create file to save log of Src MAC Address and OS of the computer commected to the network")

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/send", sendCounter)

	go capturePacket(os.Args[1], fd)

	log.Printf("Start HTTP server on localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
