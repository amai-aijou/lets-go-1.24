package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	// Instantiate a new ServeMux (the basis for our Web Server) 
	mux := http.NewServeMux()

	// Instantiate File server
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	//API Handlers for HTTP endpoints
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	return mux
}
