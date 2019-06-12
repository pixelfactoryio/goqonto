package goqonto

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	mux *http.ServeMux

	ctx = context.TODO()

	client *Client

	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = New(nil, "")
	serverUrl, _ := url.Parse(server.URL)
	client.BaseURL = serverUrl
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}
