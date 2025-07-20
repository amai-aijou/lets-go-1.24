package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
		
	// CLI flags for runtime-configurable values
	// flag.Parse() must be called *before* use of variables to store them
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Instantiate a new ServeMux (the basis for our Web Server) 
	mux := http.NewServeMux()

	// Instantiate File server
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	//API Handlers for HTTP endpoints
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

    // Info() method starting message (with listen addr as attribute)
	// flag.String (line 14) returns pointer to value, not actual value
	// pointers must be dereferenced with the * prefix. need to google this later!
    logger.Info("starting server", "addr", *addr)

	// Creates a new Web Server with ListenAndServer. seems to use "err" because
	// errors are returned through the server as non-nil entries (caight by logger.Error)
	err := http.ListenAndServe(*addr, mux)
	// Error() method logs errors returned by http.ListenAndServ; terminate with code 1
	logger.Error(err.Error())
	os.Exit(1)
}
