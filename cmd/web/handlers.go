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

	memos, err := app.memos.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

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

	// Create an instance of a `templateData` struct holding the slice of *memos*.
	data := templateData {
		Memos: memos,
	}

	// The last parameter to `Execute()` represents any dynamic data we want to pass in.
	err = ts.ExecuteTemplate(w, "base", data)
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

	// Initialize a slice containing the path to the view.tmpl.html file.
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/view.tmpl.html",
	}

	// Parse the template files...
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Create an instance of `templateData` struct, holding the *Memo* data:
	data := templateData {
		Memo: memo,
	}

	// Execute. The data `a models.Memo struct` is passed as the final parameter.
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
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
