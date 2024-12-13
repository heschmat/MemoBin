package main

import (
	"fmt"
	"net/http"
	"strconv"
)


func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"Server": "Go"}`))
	// w.Write([]byte("Welcome to MemoBin"))
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

func memoCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new memo..."))
}

func memoCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Create a new memo..."))
}
