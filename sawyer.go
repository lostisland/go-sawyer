package sawyer

import (
	"encoding/json"
	"github.com/lostisland/go-sawyer/mediatype"
	"io"
	"net/http"
	"net/url"
	"reflect"
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
	Header     http.Header
	Query      url.Values
	ErrorType  reflect.Type
}

func New(endpoint *url.URL, client *http.Client) *Client {
	if client == nil {
		client = httpClient
	}

	if len(endpoint.Path) > 0 && !strings.HasSuffix(endpoint.Path, "/") {
		endpoint.Path = endpoint.Path + "/"
	}

	return &Client{client, endpoint, make(http.Header), endpoint.Query(), nil}
}

func NewFromString(endpoint string, client *http.Client) (*Client, error) {
	e, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return New(e, client), nil
}

func (c *Client) ResolveReference(u *url.URL) *url.URL {
	absurl := c.Endpoint.ResolveReference(u)
	if len(c.Query) > 0 {
		absurl.RawQuery = mergeQueries(c.Query, absurl.Query())
	}
	return absurl
}

func (c *Client) ResolveReferenceString(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	return c.ResolveReference(u).String(), nil
}

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
