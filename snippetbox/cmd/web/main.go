package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// CLI flags for runtime-configurable values
	// flag.Parse() must be called *before* use of variables to store them

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// flag.String() returns a pointer to the flag value, not the actual value
	// pointers must be dereferenced with the * prefix. need to google this later!
	log.Print("starting server on %s", *addr) 

	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
