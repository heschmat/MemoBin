package main

import (
	"log"
	"net/http"
)


func main() {
	mux := http.NewServeMux()
	// Create a file-server which serves files out of the './ui/static' directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Register the `fileServer` as the handler for all URL paths starting with '/static/'
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Register the other application routes
	mux.HandleFunc("GET /{$}", home) // Restrict the route to exact matches on / only
	mux.HandleFunc("GET /memo/view/{id}", memoView)
	mux.HandleFunc("GET /memo/create", memoCreate)
	mux.HandleFunc("POST /memo/create", memoCreatePost)

	log.Print("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
