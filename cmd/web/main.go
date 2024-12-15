package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"text/template"

	"github.com/heschmat/MemoBin/internal/models" // {project-model-path}/internal/models

	_ "github.com/go-sql-driver/mysql" // added manually
)

// Define an application struct to hold the application-wide dependencies.
type application struct {
	logger        *slog.Logger
	memos         *models.MemoModel // `MemoModel` will be available to our handlers.
	templateCache map[string]*template.Template
}


func main() {
	// The value of the flag will be stored in the `addr` variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Define a new command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "web:changeme@/memobin?parseTime=true", "MySQL data source name")

	// If any errors are encountered during parsing, the application will be terminated.
	flag.Parse()

	// Initialize a new structured logger.
	// logger := slog.New(slog.NewTextHandler(os.Stdout, nil))  // default settings
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // min log level
		AddSource: true, // record caller location (under `source` key)
	}))

	// Pass openDB() the DSN from the cl-flag:
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// So that the connection pool is closed before the main() exits.
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}


	// Initialize a new instance of the `application` struct, containing the dependencies:
	app := &application{
		logger:        logger,
		// Initialize a `models.MemoModel` instance containing the connection pool.
		memos:         &models.MemoModel{DB: db},
		templateCache: templateCache,
	}

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

	// The value returned from the flag.String() function is a pointer to the flag value.
	logger.Info("Starting serve", "addr", *addr)
	err = http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}


// openDB() wraps sql.Open()
// and returns a sql.DB connection pool for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Verfiy everything is setup correctly.
	// db.Ping() creates a connection and checks for any error.
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
