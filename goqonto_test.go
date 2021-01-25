package goqonto

import (
	"context"
	"encoding/json"
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

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	t.Helper()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("request Body is %s, want %s", got, want)
	}
}

func testURLParseError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func testClientDefaultBaseURL(t *testing.T, c *Client) {
	t.Helper()
	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL \n got %v\n want %v\n", got, want)
	}
}

func testClientDefaultUserAgent(t *testing.T, c *Client) {
	t.Helper()
	if got, want := c.UserAgent, userAgent; got != want {
		t.Errorf("NewClient UserAgent \n got %v\n want %v\n", got, want)
	}
}

func testClientServices(t *testing.T, c *Client) {
	t.Helper()

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

// Test whether the marshaling of v produces JSON that corresponds
// to the want string.
func testJSONMarshal(t *testing.T, v interface{}, want string) {
	t.Helper()
	// Unmarshal the wanted JSON, to verify its correctness, and marshal it back
	// to sort the keys.
	u := reflect.New(reflect.TypeOf(v)).Interface()
	if err := json.Unmarshal([]byte(want), &u); err != nil {
		t.Errorf("Unable to unmarshal JSON for %v: %v", want, err)
	}
	w, err := json.Marshal(u)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", u)
	}

	// Marshal the target value.
	j, err := json.Marshal(v)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %#v", v)
	}

	if string(w) != string(j) {
		t.Errorf("json.Marshal(%q) returned %s, want %s", v, j, w)
	}
}

func testResponseMeta(t *testing.T, got *ResponseMeta, want *ResponseMeta) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ResponseMeta \n got %v\n want %v\n", got, want)
	}
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

	inBody := &TransactionsOptions{
		Slug:   "mycompany-9134",
		IBAN:   "FR761679800001000000123456",
		Status: []string{"completed"},
	}
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

func TestNewRequest_badBody(t *testing.T) {
	c := NewClient(nil)

	inURL := "v2/foo"
	badBody := make(chan struct{})

	_, err := c.NewRequest(ctx, http.MethodGet, inURL, badBody)
	if err == nil {
		t.Errorf("Expected error to be returned")
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

	want := fmt.Sprintf("%s %s", ua, userAgent)
	if got := req.Header.Get("User-Agent"); got != want {
		t.Errorf("New() UserAgent got %s\n want %s\n", got, want)
	}
}

func TestCustomBaseURL(t *testing.T) {
	baseURL := "http://localhost/foo"
	c, err := New(nil, SetBaseURL(baseURL))

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	want := baseURL
	if got := c.BaseURL.String(); got != want {
		t.Errorf("New() BaseURL got %s\n want %s\n", got, want)
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

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body got %v\n want %v\n", body, want)
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
		testMethod(t, r, http.MethodGet)
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
		t.Errorf("Completed request got %v\n want %v\n", completedReq, req)
	}

	want := `{"A":"a"}`
	if !strings.Contains(completedResp, want) {
		t.Errorf("want response to contain %v, Response = %v", want, completedResp)
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

	want := &ErrorResponse{
		Response: res,
		Message:  "m",
	}
	if !reflect.DeepEqual(err, want) {
		t.Errorf("Error got %#v\n want %#v\n", err, want)
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

	want := &ErrorResponse{
		Response: res,
	}
	if !reflect.DeepEqual(err, want) {
		t.Errorf("Error got %#v\n want %#v\n", err, want)
	}
}

func TestErrorResponse_Error(t *testing.T) {
	res := &http.Response{Request: &http.Request{}}
	err := ErrorResponse{Message: "m", Response: res}
	if err.Error() == "" {
		t.Errorf("Expected non-empty err.Error()")
	}
}
