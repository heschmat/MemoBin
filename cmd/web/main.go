package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"            // automatic form parsing
	"github.com/heschmat/MemoBin/internal/models" // {project-model-path}/internal/models

	_ "github.com/go-sql-driver/mysql" // added manually
)

// Define an application struct to hold the application-wide dependencies.
type application struct {
	debug          bool
	logger         *slog.Logger
	memos          *models.MemoModel // `MemoModel` will be available to our handlers.
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder  // holds a pointer to a `form.Decoder` instance
	sessionManager *scs.SessionManager
	users          *models.UserModel
}


func main() {
	// The value of the flag will be stored in the `addr` variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Define a new command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "web:changeme@/memobin?parseTime=true", "MySQL data source name")

	debug := flag.Bool("debug", true, "Enable debug mode")

	// If any errors are encountered during parsing, the application will be terminated.
	flag.Parse()

	// Initialize a new structured logger.
	// v1:
	// logger := slog.New(slog.NewTextHandler(os.Stdout, nil))  // default settings

	// v2:
	// logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 	Level: slog.LevelDebug, // min log level
	// 	AddSource: true, // record caller location (under `source` key)
	// }))

	// v3: use relative pass for `source`
	appDir, _ := os.Getwd()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
        AddSource: true,
        ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
            if a.Key == slog.SourceKey {
                source := a.Value.Any().(*slog.Source)
                // Convert absolute path to relative path
                relPath, err := filepath.Rel(appDir, source.File)
                if err == nil {
                    // Return a new source with relative path
                    return slog.Any(slog.SourceKey, &slog.Source{
                        File: relPath,
                        Line: source.Line,
                    })
                }
            }
            return a
        },
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


	// scs.New() returns a pointer to `SessionManager` struct
	// This holds the configuration settings for sessions.
	sessionManager := scs.New()
	// Configure `sessionManager` to use our MySQL database as the session store.
	sessionManager.Store = mysqlstore.New(db)
	// Set a lifetime of 12 hours; sessions automatically expire 12H after being created.
	sessionManager.Lifetime = 12 * time.Hour

	// Initialize a new instance of the `application` struct, containing the dependencies:
	app := &application{
		debug:          *debug,
		logger:         logger,
		// Initialize a `models.MemoModel` instance containing the connection pool.
		memos:          &models.MemoModel{DB: db},
		templateCache:  templateCache,
		// Initialize a decoder instance & add it to the application dependencies:
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
		users:          &models.UserModel{DB: db},
	}

	// Initialize a new `http.Server` struct.
	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		// Create a *log.Logger* from our standard logger handler.
		// This writes log entries at Error level.
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// The value returned from the flag.String() function is a pointer to the flag value.
	logger.Info("Starting serve", "addr", *addr)

	// Call the `.ListenAndServe()` method on the `http.Server` struct to start the server:
	err = srv.ListenAndServe()
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
