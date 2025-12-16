package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	//"strings"		// 7.3
	//"unicode/utf8"	// 7.3

	"snippetbox.nerv.com/internal/models"
	"snippetbox.nerv.com/internal/validator" // 7.5
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

	// 7.4: Initialize a new snippetCreateForm instance and pass it to the template. Great time to set default/'initial' values for the form
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

// 7.4: Define a snippetCreateForm to represent the form data and validation errors for the form fields. all struct fields
// are deliberately exported. struct fields must be exported to be read by html/template package when rendering template
// 7.5: Remove the explicity FieldErrors struct field; instead embed Validator stuct. snippetCreateForm "inherits" all fields and methods of Validator stuct
type snippetCreateForm struct {
	Title				string
	Content				string
	Expires				int
	validator.Validator
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// 7.2: First we call r.ParseForm() to add any data in POST request bodies to r.PostForm map.
	// Works the same for PUT and PATCH. If there are errors, we use app.ClientError() to send 400 to user
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// 7.2: PostForm.Get() always returns the form data as string. Since we expect a number, and want to present as integer...
	// ...we need to manually convert form data to an integer using strconv.Atoi(), and send a 400 Bad Request if that fails
	// 7.4: Get the expires value from the form as normal
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// 7.3: Initialize map to hold validation errors for the form fields (in 7.4, this is simply part of the snippetCreateForm struct)
	// 7.4: Create an instance of snippetCreateForm struct containing values from the form and an empty map for validation errors
	form := snippetCreateForm{
		Title:			r.PostForm.Get("title"),
		Content:		r.PostForm.Get("content"),
		Expires:		expires,
		// 7.5: remove this assignment:	FieldErrors:	map[string]string{},
	}

	// 7.5: since Validator struct is embedded in snippetCreateForm struct, we can call Validator.CheckField() directly on it.
	// it adds the provided key and err to FieldErrors map if check != true.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")

	// 7.5: Use Validator.Valid() to see if checks failed. If so, re-render template, passing in the form as before
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}
/*
	// 7.3: Check that title value is not blank and not more than 100 characters long. If it fails either of those checks,
	// add a message to the errors map using the field name as the key
	// 7.4: Update the validation checks so they operate on the snippetCreateForm instance
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	// 7.3: Check the 'expires' value matches one of the permitted 1/7/365 radio button values from the page
	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7, or 365"
	}

	// 7.3: If there are errors, dump them in plain-text HTTP response and return from the handler
	// 7.4: If there are validation errors, then re-render the create.tmpl template, passing in the snippetCreateForm instance
	// as dynamic data in the Form field. Use HTTP status code 422 Unprocessable Entity in response to indicate validation error
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}
*/
	// Pass data to the SnippetModel.Insert() method, receiving ID of the new record back
	// 7.4: Update line to pass the data from snippetCreateForm instance to our Insert() method.
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

