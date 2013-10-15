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

func (c *Client) ResolveReference(u *url.URL) *url.URL {
	return c.Endpoint.ResolveReference(u)
}
