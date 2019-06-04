package goqonto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/amine7536/goqonto/context"
)

// Client Qonto API Client struct
type Client struct {
	client  *http.Client
	BaseURL *url.URL

	Organizations OrganizationsService
	Transactions  TransactionsService
	Memberships   MembershipsService

	// Optional function callback
	onRequestCompleted RequestCompletionCallback
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

// Response struct
type Response struct {
	*http.Response
	Meta *ResponseMeta
}

// ResponseMeta struct
type ResponseMeta struct {
	CurrentPage int `json:"current_page,omiempty"`
	NextPage    int `json:"next_page,omiempty"`
	PrevPage    int `json:"prev_page,omiempty"`
	TotalPages  int `json:"total_pages,omiempty"`
	TotalCount  int `json:"total_count,omiempty"`
	PerPage     int `json:"per_page,omiempty"`
}

// metaRoot root key in the JSON response for meta
type metaRoot struct {
	Meta ResponseMeta `json:"meta"`
}

// Convert ResponseMeta to a string
// TODO: shouldn't Panic here
func (m ResponseMeta) String() string {
	bytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// New returns new Qonto API Client
func New(httpClient *http.Client, apiURL string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(apiURL)

	c := &Client{
		client:  httpClient,
		BaseURL: baseURL,
	}
	c.Organizations = &OrganizationsServiceOp{client: c}
	c.Transactions = &TransactionsServiceOp{client: c}
	c.Memberships = &MembershipsServiceOp{client: c}

	return c
}

// NewRequest prepare Request
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")

	return req, nil
}

func newResponse(r *http.Response) *Response {
	response := &Response{
		Response: r,
	}
	return response
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := context.DoRequestWithClient(ctx, c.client, req)
	if err != nil {
		return nil, err
	}
	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, err
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			return err
		}
	}

	return errorResponse
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}
