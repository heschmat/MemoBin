package main

import (
	"testing"
	"time"
)


func TestHumanDate(t *testing.T) {
	// Initialize a new `time.Time` object and pass it to the *humanDate* function.
	tm := time.Date(2025, 1, 10, 15, 03, 0, 0, time.UTC)
	hd := humanDate(tm)

	expectedDate := "10 Jan 2025 at 15:03"
	if hd != expectedDate {
		t.Errorf("got %q; want %q", hd, expectedDate)
	}
}
