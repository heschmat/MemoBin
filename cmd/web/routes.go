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

	// We leave the static files route unchanged.
	// Create a new middleware chain containing the middleware specific to our dynamic application routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Update the routes to use the `dynamic` middleware chain,
	// followed by the appropriate handler function.
	// N.B. the alice `ThenFunc()` method returns a `http.Handler` (rather than a `http.HandlerFunc`)
	// Registering the other application routes
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home)) // Restrict the route to exact matches on / only
	mux.Handle("GET /memo/view/{id}", dynamic.ThenFunc(app.memoView))
	mux.Handle("GET /memo/create", dynamic.ThenFunc(app.memoCreate))
	mux.Handle("POST /memo/create", dynamic.ThenFunc(app.memoCreatePost))

	// User auth routes ---------------------------------------------- //
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", dynamic.ThenFunc(app.userLogoutPost))

	// middlewares chain
	// return app.recoverPanic(app.logRequest(commonHeaders(mux)))

	// Create a middleware chain containing our *standard* middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Return the standard middleware chain followed by the *servemux*.
	return standard.Then(mux)
}
