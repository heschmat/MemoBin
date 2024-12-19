package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/heschmat/MemoBin/internal/models"
	"github.com/heschmat/MemoBin/internal/validator"
)

// N.B. struct fields must be exported (i.e., start with a capital letter)
// in order to be ready by the html/template package when rendering the template.
// This does NOT apply to maps, though.
type memoCreateForm struct {
	Title                string
	Content              string
	Expires              int
	validator.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	memos, err := app.memos.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
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

	data := app.newTemplateData(r)
	data.Memo = memo
	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) memoCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Initialize a new `memoCreateForm` instance and pass it to the template.
	// We also could set default values for the fields.
	data.Form = memoCreateForm{
		Expires: 7,
	}

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) memoCreatePost(w http.ResponseWriter, r *http.Request) {
	// Limit the request body size to 4096 bytes.
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	// if the limit above is hit, `MaxBytesreader` sets a flag on *http.ResponseWriter*
	// which instructs the server to close the underlying TCP connection.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Get the `expires` value from the *form* as normal.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := memoCreateForm {
		Title: r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 chars long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank.")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "Permitted values: 1, 7, 365")

	// If there are any validation errors,
	// re-display the `create.tmpl.html` template, passing the `memoCreateForm` instance
	// as dynamic data in the Form field.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		// To indicate that there was a validation error, we use the HTTP status code 442.
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.memos.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Redirect the user to the relevant page for the memo.
	http.Redirect(w, r, fmt.Sprintf("/memo/view/%d", id), http.StatusSeeOther)
}
