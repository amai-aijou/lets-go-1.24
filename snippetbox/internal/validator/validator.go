package validator

import (
    "slices"
    "strings"
    "unicode/utf8"
)

// 7.6: Define a new Validator struct containing a map of validation error messages for form fields
type Validator struct {
    FieldErrors map[string]string
}

// 7.6: Valid() returns true if the FieldErrors map doesn't contain any entries.
func (v *Validator) Valid() bool {
    return len(v.FieldErrors) == 0
}

// 7.6: AddFieldError() adds an error message to the FieldErrors map if no entry already exists for a given key
func (v *Validator) AddFieldError(key, message string) {
    // 7.6: Note: Need to initialize the map first, if not already initialized
    if v.FieldErrors == nil {
        v.FieldErrors = make(map[string]string)
    }

    if _, exists := v.FieldErrors[key]; !exists {
        v.FieldErrors[key] = message
    }
}

// 7.6: CheckField() adds an error message to the FieldErrors map only if a validation check is not 'ok'
func (v *Validator) CheckField(ok bool, key, message string) {
    if !ok {
        v.AddFieldError(key, message)
    }
}

// 7.6: NotBlank() returns true if a value is not an empty string
func NotBlank(value string) bool {
    return strings.TrimSpace(value) != ""
}

// 7.6: MaxChars() returns true if a value contains no more than n characters
func MaxChars(value string, n int) bool {
    return utf8.RuneCountInString(value) <= n
}

// 7.6: PermittedValue() returns true if a value is in a list of specific permitted values
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
    return slices.Contains(permittedValues, value)
}
