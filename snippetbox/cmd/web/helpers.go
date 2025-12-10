package main

import (
	"fmt" // 5.3
	"net/http"
	"runtime/debug" // 3.4: Needed for debug mode
)

// serverError helper writes an Error-level log entry, then sends 500 to user
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method	= r.Method
		uri		= r.URL.RequestURI()
		trace = string(debug.Stack()) // 3.4: Sets debug mode while using "go run" and gives stacktrace
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends 400 and similar status code and description to user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
