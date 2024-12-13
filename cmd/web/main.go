package main

import (
	"log"
	"net/http"
)


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home) // Restrict the route to exact matches on / only
	mux.HandleFunc("GET /memo/view/{id}", memoView)
	mux.HandleFunc("GET /memo/create", memoCreate)
	mux.HandleFunc("POST /memo/create", memoCreatePost)

	log.Print("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
