package main

import (
	"net/http"


	"github.com/justinas/alice" // 6.5: middleware management package	
)

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

	// 6.5: Create middleware chain containing 'standard' middleware, used for every request application receives
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// 6.2: Passes the servemux as "next" parameter to commonHeaders middleware, which later returns http.Handler
	// 6.3: Wrap the 6.2 part in logRequest(); now logRequest -> commonHeaders -> servemux -> application handler
	// 6.4: WRAP IT AGAIN BOYS!!! this time, to recover from panics
	//return app.recoverPanic(app.logRequest(commonHeaders(mux)))

	// 6.5: Instead of returning the multiple wraps, it will return the 'standard' middlewar chain followed by servemux
	return standard.Then(mux)
}
