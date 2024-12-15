package main

import (
	"path/filepath"
	"text/template"

	"github.com/heschmat/MemoBin/internal/models"
)

// Define a `templateData`
// to act as the holding structure for any dynamic data
// we want to pass to the HTML templates.
type templateData struct {
	Memo  models.Memo
	Memos []models.Memo
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
		// Create a slice containing the filepaths for our base template, any partials and the page.
		files := []string {
			"./ui/html/base.tmpl.html",
			"./ui/html/partials/nav.tmpl.html",
			page,
		}

		// Parse the files into a template set.
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		// Add the template set to the map:
		cache[name] = ts
	}

	// Return the map:
	return cache, nil
}
