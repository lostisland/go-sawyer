package sawyer

import (
	"encoding/json"
	"github.com/lostisland/go-sawyer/mediatype"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var httpClient = &http.Client{}

func init() {
	mediatype.AddDecoder("json", func(r io.Reader) mediatype.Decoder {
		return json.NewDecoder(r)
	})
	mediatype.AddEncoder("json", func(w io.Writer) mediatype.Encoder {
		return json.NewEncoder(w)
	})
}

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

func (c *Client) Do(resource interface{}, apierr interface{}, req *http.Request) *Response {
	httpres, err := c.HttpClient.Do(req)
	return buildResponse(resource, apierr, c, httpres, err)
}

func (c *Client) Get(resource interface{}, apierr interface{}, rawurl string) *Response {
	req, err := c.NewRequest("GET", rawurl, nil)
	if err != nil {
		return apiResponse(err)
	}

	return c.Do(resource, apierr, req)
}

func (c *Client) NewRequest(method string, rawurl string, body io.Reader) (*http.Request, error) {
	u, err := c.resolveReferenceString(rawurl)
	if err != nil {
		return nil, err
	}

	return http.NewRequest("GET", u, nil)
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
