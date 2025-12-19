package main

import (
	"html/template" // 5.3
	"path/filepath" // 5.3
	"time" //5.6

	"snippetbox.nerv.com/internal/models"
)

// 5.1: Define a templateData type as a holding structure for any dynamic data needed
// to be passed to the HTML templates. This is because Go's html/template package only allows
// for one piece of dynamic data to be passed in, so we wrap multiple pieces into a single struct!
type templateData struct {
	CurrentYear int			// 5.5
	Snippet		models.Snippet
	Snippets	[]models.Snippet // 5.2
	Form		any			// 7.4
	Flash		string		// 8.3
}

// 5.5: Create a humanDate function which returns a nicely-formatted string representation of time.Time value
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// 5.5: initialize a template.FuncMap value and store in a global variable. basically a string-keyed map that
// acts as a lookup table, mapping names to functions
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// 5.3: Function to parse templates at application start, then store as in-memory cache
func newTemplateCache() (map[string]*template.Template, error) {

	// 5.3: Initialize a new map to as memory cache
	cache := map[string]*template.Template{}

	// 5.3: Use filepath.glob() to get a slice of all filepaths that match pattern.
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// 5.3: Loop through the page filepaths one-by-one
	for _, page := range pages {

		// 5.3: Extract filename (like 'home.tmpl') from the full filepath and assign to name variable
		name := filepath.Base(page)

		//5.5 template.FuncMap must be registered with the template set before calling ParseFiles().
		// We have to use template.New() to create empty template set, use Funcs() to register template.FuncMap, then parse
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// 5.3: Call ParseGlob() *on this template set* to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// 5.3: Call ParseGlob() *on this template set* to add the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// 5.3: Add the template set to the map as normal...
		cache[name] = ts
	}
	return cache, nil
}
