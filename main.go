package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to MemoBin"))
}


func memoView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi((r.PathValue("id")))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	msg := fmt.Sprintf("Display a specific memo with ID %d...", id)
	w.Write([]byte(msg))
}


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home) // Restrict the route to exact matches on / only
	mux.HandleFunc("/memo/view/{id}", memoView)

	log.Print("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
