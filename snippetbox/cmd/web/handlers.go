package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"snippetbox.nerv.com/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	// Initialize a slice containing the path to the two HTML files.
	// Base template must be the *first* file!
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	// Use template.ParseFiles() to read a template in
	// use serverError() helper to log error message, then send 500 error
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err) 
		return
	}

	// Use the Execute() method to write the content as response body
	// The last parameter represents dynamic data to pass in (for now, nil)
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// 4.7: Use SnippetModel's Get() to retrieve data for record based on id (404 for no match)
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		}else {
			app.serverError(w, r, err)
		}
		return
	}

	// 4.7: Write the snippet data as plain-text HTTP response
	fmt.Fprintf(w, "%+v", snippet)

	// Old: saving in case needed?
	//fmt.Fprintf(w, "Display a specific snippet with Id %d...", id)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	// Pass data to the SnippetModel.Insert() method, receiving ID of the new record back
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

