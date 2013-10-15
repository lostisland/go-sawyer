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

func (c *Client) Get(resource interface{}, rawurl string) (*http.Response, error) {
	u, err := c.resolveReferenceString(rawurl)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
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
