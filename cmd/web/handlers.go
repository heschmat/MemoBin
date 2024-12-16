package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/heschmat/MemoBin/internal/models"
)


func (app *application) home(w http.ResponseWriter, r *http.Request) {
	memos, err := app.memos.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData()
	data.Memos = memos
	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
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

	data := app.newTemplateData()
	data.Memo = memo
	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
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
