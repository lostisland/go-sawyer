package sawyer

import (
	"github.com/lostisland/go-sawyer/hypermedia"
	"net/http"
	"net/url"
	"strings"
)

// The default httpClient used if one isn't specified.
var httpClient = &http.Client{}

// A Client wraps an *http.Client with a base url Endpoint and common header and
// query values.
type Client struct {
	HttpClient *http.Client
	Endpoint   *url.URL
	Header     http.Header
	Query      url.Values
	Cacher     Cacher
}

// New returns a new Client with a given a URL and an optional client.
func New(endpoint *url.URL, client *http.Client) *Client {
	if client == nil {
		client = httpClient
	}

	if len(endpoint.Path) > 0 && !strings.HasSuffix(endpoint.Path, "/") {
		endpoint.Path = endpoint.Path + "/"
	}

	return &Client{client, endpoint, make(http.Header), endpoint.Query(), noOpCacher}
}

// NewFromString returns a new Client given a string URL and an optional client.
func NewFromString(endpoint string, client *http.Client) (*Client, error) {
	e, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return New(e, client), nil
}

// ResolveReference resolves a URI reference to an absolute URI from an absolute
// base URI.  It also merges the query values.
func (c *Client) ResolveReference(u *url.URL) *url.URL {
	absurl := c.Endpoint.ResolveReference(u)
	if len(c.Query) > 0 {
		absurl.RawQuery = mergeQueries(c.Query, absurl.Query())
	}
	return absurl
}

// ResolveReference resolves a string URI reference to an absolute URI from an
// absolute base URI.  It also merges the query values.
func (c *Client) ResolveReferenceString(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	return c.ResolveReference(u).String(), nil
}

// Rels attempts to get the cached relations for the given request.  If it
// hasn't been cached, send a GET to the request URL, decode the response body
// to the given value, and get the relations from the value.
func (c *Client) Rels(req *Request, value interface{}) (hypermedia.Relations, *Response) {
	if rels, ok := c.Cacher.Rels(req.Request); ok {
		return rels, NewResponse(nil)
	}

	res := req.Get()
	res.Decode(value)

	return hypermedia.Rels(value), res
}

// buildRequest assembles a net/http Request using the given relative url path.
func buildRequest(c *Client, rawurl string) (*http.Request, error) {
	u, err := c.ResolveReferenceString(rawurl)
	if err != nil {
		return nil, err
	}

	httpreq, err := http.NewRequest(GetMethod, u, nil)
	for key, _ := range c.Header {
		httpreq.Header.Set(key, c.Header.Get(key))
	}
	return httpreq, err
}

// mergeQueries merges the given url.Values into a single encoded URI query
// string.
func mergeQueries(queries ...url.Values) string {
	merged := make(url.Values)
	for _, q := range queries {
		if len(q) == 0 {
			break
		}

		for key, _ := range q {
			merged.Set(key, q.Get(key))
		}
	}
	return merged.Encode()
}
