package main

import "snippetbox.nerv.com/internal/models"

// 5.1: Define a templateData type as a holding structure for any dynamic data needed
// to be passed to the HTML templates. This is because Go's html/template package only allows
// for one piece of dynamic data to be passed in, so we wrap multiple pieces into a single struct!
type templateData struct {
	Snippet models.Snippet
}
