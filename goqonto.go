// Package goqonto is a Qonto API v2 client for Go
package goqonto

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	apiVersion     = "v2"
	defaultBaseURL = "https://thirdparty.qonto.com/" + apiVersion
	userAgent      = "goqonto/" + apiVersion
	mediaType      = "application/json"
)

// Client Qonto API Client struct
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests. Defaults to the public Qonto API, but can be
	// set to a custom domain endpoint. BaseURL should
	// always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the Qonto API.
	UserAgent string

	// Reuse a single struct instead of allocating one for each service on the heap.
	common service

	// Services used for talking to different parts of the Qonto API.
	Organizations *OrganizationsService
	Transactions  *TransactionsService
	Memberships   *MembershipsService
	Attachments   *AttachmentsService
	Labels        *LabelsService

	// Optional function callback
	onRequestCompleted RequestCompletionCallback
}

type service struct {
	client *Client
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
	CurrentPage int `json:"current_page,omitempty"`
	NextPage    int `json:"next_page,omitempty"`
	PrevPage    int `json:"prev_page,omitempty"`
	TotalPages  int `json:"total_pages,omitempty"`
	TotalCount  int `json:"total_count,omitempty"`
	PerPage     int `json:"per_page,omitempty"`
}

// metaRoot root key in the JSON response for meta
type metaRoot struct {
	Meta ResponseMeta `json:"meta"`
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewClient returns new Qonto API Client
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}
	c.common.client = c

	c.Organizations = (*OrganizationsService)(&c.common)
	c.Transactions = (*TransactionsService)(&c.common)
	c.Memberships = (*MembershipsService)(&c.common)
	c.Attachments = (*AttachmentsService)(&c.common)
	c.Labels = (*LabelsService)(&c.common)

	return c
}

// ClientOpt are options for New.
type ClientOpt func(*Client) error

// New returns a new Qonto ThirdParty API client instance.
func New(httpClient *http.Client, opts ...ClientOpt) (*Client, error) {
	c := NewClient(httpClient)
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// SetBaseURL is a client option for setting the base URL.
func SetBaseURL(bu string) ClientOpt {
	return func(c *Client) error {
		u, err := url.Parse(bu)
		if err != nil {
			return err
		}

		c.BaseURL = u
		return nil
	}
}

// SetUserAgent is a client option for setting the user agent.
func SetUserAgent(ua string) ClientOpt {
	return func(c *Client) error {
		c.UserAgent = fmt.Sprintf("%s %s", ua, c.UserAgent)
		return nil
	}
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

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)

	return req, nil
}

// OnRequestCompleted sets request completion callback
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
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
	resp, err := DoRequestWithClient(ctx, c.client, req)
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

// DoRequest submits an HTTP request.
func DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	return DoRequestWithClient(ctx, http.DefaultClient, req)
}

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(
	ctx context.Context,
	client *http.Client,
	req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
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
			errorResponse.Message = string(data)
			errorResponse.Code = errorResponse.Response.StatusCode
		}
	}

	return errorResponse
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v", e.Response.Request.Method, e.Response.Request.URL, e.Response.StatusCode, e.Message)
}
