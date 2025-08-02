package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

  // Application struct to hold app-wide dependencies
type application struct {
	logger *slog.Logger
}

func main() {
		
	// CLI flags for runtime-configurable values
	// flag.Parse() must be called *before* use of variables to store them
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Create DSN (Data Source Name) for Go MySQL driver
	dsn := flag.String("dsn", "web:example-pass@/snippetbox?parseTime=true", "MySQL data source name")
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

	// Instantiate a new application struct containinng all dependencies
	app := &application{
		logger: logger,
	}

    // Info() method starting message (with listen addr as attribute)
	// flag.String (line 14) returns pointer to value, not actual value
	// pointers must be dereferenced with the * prefix. need to google this later!
    logger.Info("starting server", "addr", *addr)

	// Creates a new Web Server with ListenAndServer. seems to use "err" because
	// errors are returned through the server as non-nil entries (caight by logger.Error)
	err = http.ListenAndServe(*addr, app.routes())

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
