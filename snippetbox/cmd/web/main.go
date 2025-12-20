package main

import (
	"database/sql"
	"flag"
	"html/template" // 5.3
	"log/slog"
	"net/http"
	"os"
	"time" // 8.2

	// Import the models package created in internal/models
	"snippetbox.nerv.com/internal/models"

	
	"github.com/alexedwards/scs/mysqlstore"	// 8.2
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"	// 7.6
	_ "github.com/go-sql-driver/mysql"
)

  // Application struct to hold app-wide dependencies
type application struct {
	logger			*slog.Logger
	snippets		*models.SnippetModel
	templateCache	map[string]*template.Template
	formDecoder		*form.Decoder
	sessionManager	*scs.SessionManager		// 8.2
}

func main() {
	// CLI flags for runtime-configurable values
	// flag.Parse() must be called *before* use of variables to store them
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Create DSN (Data Source Name) for Go MySQL driver
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initiate Database connection pool and DB driver for Go
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Close the DB connection pool (before the main function exits)
	defer db.Close()

	// 5.3: Initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// 7.6: Initialize a decoder instance
	formDecoder := form.NewDecoder()

	// 8.2: initialize a new session manager, then configure it to use MySQL DB as session store.
	// lifetime: 12hrs (expires 12hrs after being created)
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Instantiate a new application struct containing all dependencies
	// AND: Instantiate a new SnippetModel instance with connection pool
	app := &application{
		logger:			logger,
		snippets:		&models.SnippetModel{DB: db},
		templateCache:	templateCache,
		formDecoder:	formDecoder,
		sessionManager:	sessionManager, // 8.2
	}

	// 9.1: Manually create http.Server struct for more control over server than http.ListenAndServe() can give.
	// Initialize http.Server struct. Set Addr and Handler fields to use flags from above
	srv := &http.Server{
		Addr:		*addr,
		Handler:	app.routes(),
	}

    // Info() method starting message (with listen addr as attribute)
	// flag.String (line 14) returns pointer to value, not actual value
	// pointers must be dereferenced with the * prefix. need to google this later!
    logger.Info("starting server", "addr", *addr)

	// Creates a new Web Server with ListenAndServer. seems to use "err" because
	// errors are returned through the server as non-nil entries (caight by logger.Error)
	// 9.1: Update from http.ListenAndServe to new srv.ListenAndServe(), with custom struct above
	err = srv.ListenAndServe()

	// Error() method logs errors returned by http.ListenAndServ; terminate with code 1
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
