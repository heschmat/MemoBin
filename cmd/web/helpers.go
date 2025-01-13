package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
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

	statusCode := http.StatusInternalServerError
	if app.debug {
		body := fmt.Sprintf("%s\n%s", err, trace)
		http.Error(w, body, statusCode)
		return
	}

	http.Error(w, http.StatusText(statusCode), statusCode)
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
		err := fmt.Errorf("the template *%s* does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	// First write the template to the buffer.
	// If there's an error, call `serverError()` and return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear:     time.Now().Year(),
		// Add the flash message to the template data, if one exists.
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken: nosurf.Token(r),
	}
}
