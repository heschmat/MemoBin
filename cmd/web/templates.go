package main

import "github.com/heschmat/MemoBin/internal/models"

// Define a `templateData`
// to act as the holding structure for any dynamic data
// we want to pass to the HTML templates.
type templateData struct {
	Memo  models.Memo
	Memos []models.Memo
}
