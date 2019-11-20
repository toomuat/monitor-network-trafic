package main

import (
	"log"
	"net/http"
)

func main() {
	// http.Handle("/", http.FileServer(http.Dir("./static")))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.HandleFunc("/chart", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/cht.html")
	})

	log.Printf("Start HTTP server")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
