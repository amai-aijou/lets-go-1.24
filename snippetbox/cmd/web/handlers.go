package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.nerv.com/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// 6.4: BREAK SOMETHING!!! THIS IS A PURPOSEFUL ERROR FOR LEARNING PURPOSES
	// panic("oops! something went wrong...very, very wrong") //Deliberate panic

	// 4.8: Use SnippetModel.Latest() method, dumping snippet contents to HTTP response
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// 5.5: Call newTemplateData() helper to get templateData struct containing 'default' data
	// (right now just the current year), and add the snippets slice to it
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// 5.5: Pass the data to the render() helper as normal
	app.render(w, r, http.StatusOK, "home.tmpl", data)
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

	// 5.5: call newTemplateData() helper to get default data, and addd snppets slice to it
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", data)
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

