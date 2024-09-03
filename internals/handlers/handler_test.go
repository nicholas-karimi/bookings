package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	Key   string
	Value string
}

var tests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"generals", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"reservations", "/make-reservations", "GET", []postData{}, http.StatusOK},
	{"summary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	{"post-search-availability", "/search-availability", "POST", []postData{
		{Key: "start", Value: "2024-09-03"},
		{Key: "end", Value: "2024-09-04"},
	}, http.StatusOK},
	{"post-search-availability-json", "/search-avalibility-json", "POST", []postData{
		{Key: "start", Value: "2024-09-03"},
		{Key: "end", Value: "2024-09-04"},
	}, http.StatusOK},
	{"post-make-reservations", "/make-reservations", "POST", []postData{
		{Key: "first_name", Value: "John"},
		{Key: "last_name", Value: "Doe"},
		{Key: "email", Value: "j@j.com"},
		{Key: "phone", Value: "123456789"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {

	routes := getRoutes()
	test_client := httptest.NewTLSServer(routes)
	defer test_client.Close()

	for _, e := range tests {
		if e.method == "GET" {
			resp, err := test_client.Client().Get(test_client.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected status code %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}

			for _, x := range e.params {
				values.Add(x.Key, x.Value)
			}

			resp, err := test_client.Client().PostForm(test_client.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected status code %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}

}
