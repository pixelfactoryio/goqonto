package goqonto

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var (
	mux    *http.ServeMux
	ctx    = context.TODO()
	client *Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(nil)
	serverURL, _ := url.Parse(server.URL)
	client.BaseURL = serverURL
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if got, want := expected, r.Method; got != want {
		t.Errorf("Request Method \n got %v\n want %v\n", got, want)
	}
}

func testResponseMeta(t *testing.T, got *ResponseMeta, want *ResponseMeta) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ResponseMeta \n got %v\n want %v\n", got, want)
	}
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func testClientDefaultBaseURL(t *testing.T, c *Client) {
	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL \n got %v\n want %v\n", got, want)
	}
}

func testClientDefaultUserAgent(t *testing.T, c *Client) {
	if got, want := c.UserAgent, userAgent; got != want {
		t.Errorf("NewClient UserAgent \n got %v\n want %v\n", got, want)
	}
}

func testClientServices(t *testing.T, c *Client) {
	services := []string{
		"Organizations",
		"Transactions",
		"Memberships",
		"Attachments",
	}

	cp := reflect.ValueOf(c)
	cv := reflect.Indirect(cp)

	for _, s := range services {
		if cv.FieldByName(s).IsNil() {
			t.Errorf("c.%s shouldn't be nil", s)
		}
	}

}

func testClientDefaults(t *testing.T, c *Client) {
	testClientDefaultBaseURL(t, c)
	testClientDefaultUserAgent(t, c)
	testClientServices(t, c)
}

func TestNew(t *testing.T) {
	c, err := New(nil)

	if err != nil {
		t.Fatalf("New(): %v", err)
	}

	testClientDefaults(t, c)
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	c2 := NewClient(nil)
	if c.client == c2.client {
		t.Error("NewClient returned same http.Clients, but they should differ")
	}

	testClientDefaults(t, c)
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL := "v2/foo"
	outURL := defaultBaseURL + "/foo"

	inBody := &TransactionsOptions{Slug: "mycompany-9134", IBAN: "FR761679800001000000123456", Status: []string{"completed"}}
	outBody := `{"slug":"mycompany-9134","iban":"FR761679800001000000123456","status":["completed"]}` + "\n"

	req, _ := c.NewRequest(ctx, http.MethodGet, inURL, inBody)

	// test relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%v) \n got %v\n want %v\n", inURL, got, want)
	}

	// test body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%v) \n got %v\n want %v\n", inBody, got, want)
	}

	// test default user-agent is attached to the request
	userAgent := req.Header.Get("User-Agent")
	if got, want := c.UserAgent, userAgent; got != want {
		t.Errorf("NewRequest() \n got %v\n want %v\n", got, want)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient(nil)
	_, err := c.NewRequest(ctx, http.MethodGet, ":", nil)
	testURLParseError(t, err)
}

func TestNewRequest_withCustomUserAgent(t *testing.T) {
	ua := "testing/0.0.1"
	c, err := New(nil, SetUserAgent(ua))

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	req, _ := c.NewRequest(ctx, http.MethodGet, "/foo", nil)

	expected := fmt.Sprintf("%s %s", ua, userAgent)
	if got := req.Header.Get("User-Agent"); got != expected {
		t.Errorf("New() UserAgent = %s; expected %s", got, expected)
	}
}

func TestCustomBaseURL(t *testing.T) {
	baseURL := "http://localhost/foo"
	c, err := New(nil, SetBaseURL(baseURL))

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	expected := baseURL
	if got := c.BaseURL.String(); got != expected {
		t.Errorf("New() BaseURL = %s; expected %s", got, expected)
	}
}

func TestCustomBaseURL_badURL(t *testing.T) {
	baseURL := ":"
	_, err := New(nil, SetBaseURL(baseURL))

	testURLParseError(t, err)
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	body := new(foo)
	_, err := client.Do(context.Background(), req, body)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}

	expected := &foo{"a"}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Response body = %v, expected %v", body, expected)
	}
}

func TestDo_httpError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

// Test handling of an error caused by the internal http client's Do()
// function.
func TestDo_redirectLoop(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected a URL error; got %#v.", err)
	}
}

func TestDo_completion_callback(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	body := new(foo)
	var completedReq *http.Request
	var completedResp string
	client.OnRequestCompleted(func(req *http.Request, resp *http.Response) {
		completedReq = req
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			t.Errorf("Failed to dump response: %s", err)
		}
		completedResp = string(b)
	})
	_, err := client.Do(context.Background(), req, body)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}
	if !reflect.DeepEqual(req, completedReq) {
		t.Errorf("Completed request = %v, expected %v", completedReq, req)
	}
	expected := `{"A":"a"}`
	if !strings.Contains(completedResp, expected) {
		t.Errorf("expected response to contain %v, Response = %v", expected, completedResp)
	}
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body: ioutil.NopCloser(strings.NewReader(`{"message":"m",
			"errors": [{"resource": "r", "field": "f", "code": "c"}]}`)),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Fatalf("Expected error response.")
	}

	expected := &ErrorResponse{
		Response: res,
		Message:  "m",
	}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("Error = %#v, expected %#v", err, expected)
	}
}

// ensure that we properly handle API errors that do not contain a response
// body
func TestCheckResponse_noBody(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Errorf("Expected error response.")
	}

	expected := &ErrorResponse{
		Response: res,
	}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("Error = %#v, expected %#v", err, expected)
	}
}

func TestErrorResponse_Error(t *testing.T) {
	res := &http.Response{Request: &http.Request{}}
	err := ErrorResponse{Message: "m", Response: res}
	if err.Error() == "" {
		t.Errorf("Expected non-empty ErrorResponse.Error()")
	}
}
