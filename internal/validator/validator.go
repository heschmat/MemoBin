package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

// Define a map of validation error messages for the form fields.
// *NonFieldErrors* holds any validation errors, not related to a specific form field.
type Validator struct {
	FieldErrors    map[string]string
	NonFieldErrors []string
}


var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Retrun true if no error - field-specific or not - is registered.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
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

func (v * Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
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

func MinChars(s string, n int) bool {
	return utf8.RuneCountInString(s) >= n
}

func Matches(s string, rx *regexp.Regexp) bool {
	return rx.MatchString(s)
}
