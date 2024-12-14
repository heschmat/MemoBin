package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/heschmat/MemoBin/internal/models"
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
		// app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.serverError(w, r, err)
		return
	}

	// The last parameter to `Execute()` represents any dynamic data we want to pass in.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		// app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.serverError(w, r, err)
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

	memo, err := app.memos.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			// If no matching record is found, return a 404 Not Found response.
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Write the *memo* data as a plain-text HTTP response body.
	fmt.Fprintf(w, "%+v", memo)
}

func (app *application) memoCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new memo..."))
}

func (app *application) memoCreatePost(w http.ResponseWriter, r *http.Request) {
	// Dummy data for now.
	title := "TypeScript"
	content := "JS with types, enhancing safety for large-scale applications."
	expires := 1
	id, err := app.memos.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect the user to the relevant page for the memo.
	http.Redirect(w, r, fmt.Sprintf("/memo/view/%d", id), http.StatusSeeOther)
}