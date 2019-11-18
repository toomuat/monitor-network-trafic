package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// http.HandleFunc("/chart", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./static/cht.html")
	// 	http.ServeFile(w, r, "./static/cht.css")
	// 	http.ServeFile(w, r, "./static/reset.css")
	// 	http.ServeFile(w, r, "./static/cht.js")
	// 	http.ServeFile(w, r, "./static/chartjs-plugin-streaming.min.js")
	// })

	// http.HandleFunc("/index", func(w http.ResponseWriter, req *http.Request) {
	// 	w.Write([]byte("Hello world"))
	// })

	log.Printf("Start HTTP server")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

	// err := http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// }))

	// if err != nil {
	// 	panic(err)
	// }
}
