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

func (c *Client) Do(resource interface{}, apierr interface{}, req *http.Request) (*http.Response, error) {
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	return res, c.decode(resource, apierr, res)
}

func (c *Client) Get(resource interface{}, apierr interface{}, rawurl string) (*http.Response, error) {
	req, err := c.NewRequest("GET", rawurl, nil)
	if err != nil {
		return nil, err
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

func (c *Client) decode(resource interface{}, apierr interface{}, res *http.Response) error {
	if UseApiError(res.StatusCode) {
		return c.decodeResource(apierr, res)
	}
	return c.decodeResource(resource, res)
}

func (c *Client) decodeResource(resource interface{}, res *http.Response) error {
	if resource == nil {
		return nil
	}

	dec, err := c.decoder(res)
	if err != nil {
		return err
	} else if dec != nil {
		return dec.Decode(resource)
	}

	return nil
}

func (c *Client) decoder(res *http.Response) (Decoder, error) {
	mt, err := c.mediaType(res)
	if err != nil {
		return nil, err
	}

	if decfunc, ok := c.Decoders[mt.Format]; ok {
		return decfunc(res.Body), nil
	}
	return nil, nil
}

func (c *Client) mediaType(res *http.Response) (*mediatype.MediaType, error) {
	if ctype := res.Header.Get("Content-Type"); len(ctype) > 0 {
		return mediatype.Parse(ctype)
	}
	return nil, nil
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
