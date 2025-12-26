package validator

import (
	"regexp"	// 10.3
    "slices"
    "strings"
    "unicode/utf8"
)

// 10.3: Function parses RegEx pattern to sanity-check email address format. Returns a pointer to a 'compiled' regexp.Regexp type (or panics!).
// Note: Parses pattern once at startup and storing the compiled 8regexp.Regexp in a variable is more performant than re-parsing each time it's needed
// Note: Regex is the pattern currently recommended by W3C and WHATWG for email validation.
// Note: Pattern is an interpreted string literal, so it needs double-escape chars \\. Raw string literal uses `` meaning it can't handle that character (notice no " in the pattern!)
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// 10.3: MinChars() returns true is a char contains at least n characters
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// 10.3: Matches() returns true if a value matches with a provided compiled regex pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// 7.6: Define a new Validator struct containing a map of validation error messages for form fields
type Validator struct {
	NonFieldErrors	[]string // 10.4
    FieldErrors		map[string]string 
}

// 7.6: Valid() returns true if the FieldErrors map doesn't contain any entries.
// 10.4: Add NonFieldErrors. Now both must return zero, meaning there are no errors of either type
func (v *Validator) Valid() bool {
    return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// 10.4: AddNonFieldError() helper for any error not related to a field in the form
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
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
