package main

import (
	"bytes" // 5.4
	"errors" // 7.6
	"fmt" // 5.3
	"net/http"
	"time" // 5.5
	"runtime/debug" // 3.4: Needed for debug mode

	"github.com/go-playground/form/v4" // 7.6
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

// 5.5: returns a templateData struct initialized with the current year
func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData {
		CurrentYear: time.Now().Year(),
	}
}

// 7.6: decodePostForm() helper method for go-playground/form. The second parameter, dst, is the target destination to decode form data to
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// 7.6: call ParseForm() on the request
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// 7.6: Call Decode(), passing target destination as the first parameter
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// 7.6: If we use an invalid target dst, return an error with type form.InvalidDecoderError. We use errors.As() to check for this and panic
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// 7.6: For all other errors, return them as normal
		return err
	}

	return nil
}
