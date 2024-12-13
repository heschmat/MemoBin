package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)


func main() {
	// The value of the flag will be stored in the `addr` variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")
	// If any errors are encountered during parsing, the application will be terminated.
	flag.Parse()

	// Initialize a new structured logger.
	// logger := slog.New(slog.NewTextHandler(os.Stdout, nil))  // default settings
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // min log level
		AddSource: true, // record caller location (under `source` key)
	}))

	mux := http.NewServeMux()
	// Create a file-server which serves files out of the './ui/static' directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Register the `fileServer` as the handler for all URL paths starting with '/static/'
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Register the other application routes
	mux.HandleFunc("GET /{$}", home) // Restrict the route to exact matches on / only
	mux.HandleFunc("GET /memo/view/{id}", memoView)
	mux.HandleFunc("GET /memo/create", memoCreate)
	mux.HandleFunc("POST /memo/create", memoCreatePost)

	// The value returned from the flag.String() function is a pointer to the flag value.
	logger.Info("Starting serve", "addr", *addr)
	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
