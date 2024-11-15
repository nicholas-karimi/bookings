package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {

	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/some-url", nil)

	r.PostForm = postedData
	form = New(r.PostForm)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows not valid when required fields present")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)

	// form.Has("a", r)
	// if !form.Valid() {
	// 	t.Error("form shows valid when 'a' is present")
	// }
	has := form.Has("a", r)

	if has {
		t.Error("form shows 'a' when it should not be present")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	has = form.Has("a", r)

	if !has {
		t.Error("form shows 'a' when it should be present")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)

	form.MinLength("x", 10, r)

	if form.Valid() {
		t.Error("form shows min length for non-existent filed")
	}

	postedValues := url.Values{}
	postedValues.Add("some_field", "some value")

	form = New(postedValues)

	form.MinLength("some_field", 100, r)

	if form.Valid() {
		t.Error("shows min-lent of 100 met when data is shorter.")
	}

	postedValues = url.Values{}
	postedValues.Add("another_field", "abc1234")

	form = New(postedValues)
	form.MinLength("another_field", 1, r)

	if !form.Valid() {
		t.Error("shows min length of 1 is not met when it is")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)

	form.IsEmail("x", r)
	if form.Valid() {
		t.Error("form shows valid email for non existent field")
	}

	postedValues := url.Values{}
	postedValues.Add("email", "nkarimi@goland.com")

	form = New(postedValues)

	form.IsEmail("email", r)

	if !form.Valid() {
		t.Error("got an invalid email when we should have not")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "n")

	form = New(postedValues)

	form.IsEmail("email", r)

	if form.Valid() {
		t.Error("got valid for invalid email address")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("shows have an error, but did not get one.")
	}
}
