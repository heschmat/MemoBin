package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

// Define a map of validation error messages for the form fields.
type Validator struct {
	FieldErrors map[string]string
}

// Retrun true if the FieldErrors map doesn't contain any entries
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// So long as no entry already exists for the given key, add error msg.
func (v *Validator) AddFieldError(key, msg string) {
	// Initialize the map first, if it's not already initialized.
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = msg
	}
}

func (v *Validator) CheckField(ok bool, key, msg string) {
	if !ok {
		v.AddFieldError(key, msg)
	}
}

func NotBlank(val string) bool {
	return strings.TrimSpace(val) != ""
}

func MaxChars(val string, n int) bool {
	return utf8.RuneCountInString(val) <= n
}

func PermittedValue[T comparable](val T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, val)
}
