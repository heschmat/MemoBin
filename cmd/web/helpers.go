package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// `serverError` helper writes a log entry at Error level
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError (w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		// The stack trace returns a byte slice => convert to string
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}


// `clientError` helper sends a specific status code & corresponding description to the user.
// e.g., where there's a problem with the user's request, send responses like `400 "Bad Request`
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}


func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// Retrieve the appropriate template set from the cache base on the page name.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("The template *%s* does not exist.", page)
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
