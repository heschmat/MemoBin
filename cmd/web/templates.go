package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/heschmat/MemoBin/internal/models"
)

// Define a `templateData`
// to act as the holding structure for any dynamic data
// we want to pass to the HTML templates.
type templateData struct {
	CurrentYear int
	Memo        models.Memo
	Memos       []models.Memo
	Form        any
	Flash       string
	IsAuthenticated bool
	CSRFToken    string
}

// YYYY-MM-DD HH:MM:SS +0000 UTC => 16 Dec 2024 at 12:21
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}


func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Get a slice of all filepaths that match the pattern: "./ui/html/pages/*.tmpl.html"
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the filename from the fullpath
		name := filepath.Base(page)

		// `template.FuncMap` must be registered with the template set before calling `.ParseFiles()`
		// hence, 1) create an empty template set via `template.New()`
		// 2) register the `template.FuncMap` via `.Funcs()`
		// 3) parse the file as normal.
		// Create a slice containing the filepaths for our base template, any partials and the page.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map:
		cache[name] = ts
	}

	// Return the map:
	return cache, nil
}
