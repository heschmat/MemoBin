package main

import (
	"net/http"

	"github.com/justinas/alice"
)


func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	// Create a file-server which serves files out of the './ui/static' directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Register the `fileServer` as the handler for all URL paths starting with '/static/'
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Register the other application routes
	mux.HandleFunc("GET /{$}", app.home) // Restrict the route to exact matches on / only
	mux.HandleFunc("GET /memo/view/{id}", app.memoView)
	mux.HandleFunc("GET /memo/create", app.memoCreate)
	mux.HandleFunc("POST /memo/create", app.memoCreatePost)

	// middlewares chain
	// return app.recoverPanic(app.logRequest(commonHeaders(mux)))

	// Create a middleware chain containing our *standard* middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Return the standard middleware chain followed by the *servemux*.
	return standard.Then(mux)
}
