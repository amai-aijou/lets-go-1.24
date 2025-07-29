package main

import (
	"net/http"
	"runtime/debug"
)

// serverError helper writes an Error-level log entry, then sends 500 to user
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method	= r.Method
		uri		= r.URL.RequestURI()
		trace = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends 400 and similar status code and description to user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
