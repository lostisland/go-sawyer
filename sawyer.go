package sawyer

import (
	"net/http"
	"net/url"
	"strings"
)

var httpClient = &http.Client{}

type Client struct {
	HttpClient *http.Client
	Endpoint   *url.URL
}

func New(endpoint *url.URL, client *http.Client) *Client {
	if client == nil {
		client = httpClient
	}

	if len(endpoint.Path) > 0 && !strings.HasSuffix(endpoint.Path, "/") {
		endpoint.Path = endpoint.Path + "/"
	}

	return &Client{client, endpoint}
}

func NewFromString(endpoint string, client *http.Client) (*Client, error) {
	e, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return New(e, client), nil
}

func (c *Client) Get(resource interface{}, apierror interface{}, rawurl string) *Response {
	u, err := c.resolveReferenceString(rawurl)
	if err != nil {
		return &Response{ResponseError: err}
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return &Response{ResponseError: err}
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return &Response{ResponseError: err}
	}

	return &Response{Response: res}
}

func (c *Client) ResolveReference(u *url.URL) *url.URL {
	return c.Endpoint.ResolveReference(u)
}

func (c *Client) resolveReferenceString(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	return c.ResolveReference(u).String(), nil
}

type Response struct {
	ResponseError error
	*http.Response
}

func (r *Response) IsError() bool {
	if r.ResponseError != nil {
		return true
	}
	return false
}

func (r *Response) Error() string {
	if r.IsError() {
		return r.ResponseError.Error()
	}
	return ""
}
