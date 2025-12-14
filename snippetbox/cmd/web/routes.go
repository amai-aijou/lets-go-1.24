package main

import "net/http"

// 6.2: Wrap servemux with commonHeaders middleware to add Headers to each request.
// Now returns http.Handler instead of *http.ServeMux
func (app *application) routes() http.Handler {
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

	// 6.2: Passes the servemux as "next" parameter to commonHeaders middleware, which later returns http.Handler
	// 6.3: Wrap the 6.2 part in logRequest(); now logRequest -> commonHeaders -> servemux -> application handler
	return app.logRequest(commonHeaders(mux))
}
