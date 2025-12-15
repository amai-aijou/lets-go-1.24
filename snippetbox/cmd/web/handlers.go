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
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// 7.2: First we call r.ParseForm() to add any data in POST request bodies to r.PostForm map.
	// Works the same for PUT and PATCH. If there are errors, we use app.ClientError() to send 400 to user
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// 7.2: Use r.PostForm.Get() to retrieve title and content from r.PostForm map
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// 7.2: PostForm.Get() always returns the form data as string. Since we expect a number, and want to present as integer...
	// ...we need to manually convert form data to an integer using strconv.Atoi(), and send a 400 Bad Request if that fails
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Pass data to the SnippetModel.Insert() method, receiving ID of the new record back
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

