package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form custom struct and embed url.Values
type Form struct {
	url.Values
	Errors errors
}

// Valid rerturns true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0

}

// Required check for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank.")
		}
	}
}

// New initialize form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		// f.Errors.Add(field, "This field cannot be blank.")
		// -> use to check and return if form has particular field eg checkbox which is returned only when checked
		return false
	}

	return true
}

// MinLength check for string minimum length
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)

	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long.", length))
		return false
	}
	return true
}

// IsEmail check for valid email address
func (f *Form) IsEmail(field string, r *http.Request) {

	if !govalidator.IsEmail(r.Form.Get(field)) {
		f.Errors.Add(field, "Invalid email address.")
	}
}
