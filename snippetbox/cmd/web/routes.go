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

	// 8.2: New middleware chain for dynamic application routes, like Session Manager, since it doesn't apply to all routes
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	//API Handlers for HTTP endpoints
	// 8.2: Update to use dynamic middleware chain. Since alice.ThenFunc returns an http.Handler (instead of http.HandlerFunc),
	// we need to register the routes using mux.Handle() instead of mux.HandleFunc()
	mux.Handle("GET /{$}",					dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}",	dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup",			dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup",			dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login",			dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login",			dynamic.ThenFunc(app.userLoginPost))


	// 10.6: Protected (authenticated-only) application routes, with a new "protected" middleware chain via the requireAuthentication middleware
	protected := dynamic.Append(app.requireAuthentication)

	// 10.1: New Routes for User Authentication and management
	mux.Handle("GET /snippet/create",		protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create",		protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout",			protected.ThenFunc(app.userLogoutPost))

	// 6.5: Create middleware chain containing 'standard' middleware, used for every request application receives
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// 6.2: Passes the servemux as "next" parameter to commonHeaders middleware, which later returns http.Handler
	// 6.3: Wrap the 6.2 part in logRequest(); now logRequest -> commonHeaders -> servemux -> application handler
	// 6.4: WRAP IT AGAIN BOYS!!! this time, to recover from panics
	//return app.recoverPanic(app.logRequest(commonHeaders(mux)))

	// 6.5: Instead of returning the multiple wraps, it will return the 'standard' middlewar chain followed by servemux
	return standard.Then(mux)
}
