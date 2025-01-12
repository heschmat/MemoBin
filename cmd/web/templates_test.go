package main

import (
	"testing"
	"time"
)

// func TestHumanDate(t *testing.T) {
// 	// Initialize a new `time.Time` object.
// 	tm := time.Date(2025, 1, 10, 15, 03, 0, 0, time.UTC)
// 	hd := humanDate(tm)

// 	expectedDate := "10 Jan 2025 at 15:03"
// 	if hd != expectedDate {
// 		t.Errorf("got %q; want %q", hd, expectedDate)
// 	}
// }

func TestHumanDate(t *testing.T) {
	// Create a slice of anonymous structs containing the test case name,
	// input to our humanDate() function (the tm field)
	// and expected output (the want field)
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2025, 1, 10, 15, 03, 0, 0, time.UTC),
			want: "10 Jan 2025 at 15:03",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2025, 1, 10, 15, 03, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "10 Jan 2025 at 14:03", // CET is one hour ahead of UTC
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			if hd != tt.want {
				t.Errorf("got %q; want %q", hd, tt.want)
			}
		})
	}
}

// go test -v ./cmd/web
