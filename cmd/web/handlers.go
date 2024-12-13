package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)


func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	// Initialize a slice containing the paths to the two files:
	files := []string {
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// The last parameter to `Execute()` represents any dynamic data we want to pass in.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(`{"Server": "Go"}`))
	// w.Write([]byte("Welcome to MemoBin"))
}


func (app *application) memoView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi((r.PathValue("id")))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	msg := fmt.Sprintf("Display a specific memo with ID %d...", id)
	w.Write([]byte(msg))
}

func (app *application) memoCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new memo..."))
}

func (app *application) memoCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Create a new memo..."))
}
