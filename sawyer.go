package sawyer

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var httpClient = &http.Client{}

type Client struct {
	HttpClient *http.Client
	Endpoint   *url.URL
	Decoders   map[string]DecoderFunc
}

type DecoderFunc func(r io.Reader) Decoder

type Decoder interface {
	Decode(v interface{}) error
}

func New(endpoint *url.URL, client *http.Client) *Client {
	if client == nil {
		client = httpClient
	}

	if len(endpoint.Path) > 0 && !strings.HasSuffix(endpoint.Path, "/") {
		endpoint.Path = endpoint.Path + "/"
	}

	decoders := map[string]DecoderFunc{
		"json": func(r io.Reader) Decoder {
			return json.NewDecoder(r)
		},
	}
	return &Client{client, endpoint, decoders}
}

func NewFromString(endpoint string, client *http.Client) (*Client, error) {
	e, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return New(e, client), nil
}

func (c *Client) Get(resource interface{}, apierr interface{}, rawurl string) (*http.Response, error) {
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
		return res, err
	}
	defer res.Body.Close()

	return res, c.decode(resource, apierr, res)
}

func (c *Client) ResolveReference(u *url.URL) *url.URL {
	return c.Endpoint.ResolveReference(u)
}

func (c *Client) decode(resource interface{}, apierr interface{}, res *http.Response) error {
	// TODO: content type negotiation to find the right decoder
	dec := c.Decoders["json"](res.Body)

	if UseApiError(res.StatusCode) {
		return dec.Decode(apierr)
	}

	return dec.Decode(resource)
}

func (c *Client) resolveReferenceString(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	return c.ResolveReference(u).String(), nil
}

func UseApiError(status int) bool {
	switch {
	case status > 199 && status < 300:
		return false
	case status == 304:
		return false
	case status == 0:
		return false
	}
	return true
}
