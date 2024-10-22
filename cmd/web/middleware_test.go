package main

import (
	"net/http"
	"testing"
)

func TestNosurf(t *testing.T) {

	var myHandler myHandler
	handler := NoSurf(&myHandler)

	switch v := handler.(type) {
	case http.Handler:
		// do mothing

	default:
		t.Errorf("type is not http.Handler, type is %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myHandler myHandler

	handler := SessionLoad(&myHandler)

	switch ty := handler.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("type is not http.Handler, type is %T", ty)
	}
}
