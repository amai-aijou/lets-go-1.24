package main

import (
	"bytes" // 5.4
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

	// 5.4: Initialize new buffer
	buf := new(bytes.Buffer)

	// 5.4: Write the template to buffer, instead of straight to http.ResponseWriter.
	// If error, call serverError() helper, then return
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// 5.4: If template writes to buffer without error, it's safe to write HTTP status code to http.ResponseWriter!
	w.WriteHeader(status)

	// 5.4: Write contents of buffer to http.ResponseWriter.
	// NOTE: This is another time where we pass http.ResponseWriter to a function that takes an io.Writer.
	buf.WriteTo(w)
}
