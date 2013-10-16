package sawyer

import (
	"net/http"
)

type Request struct {
	Client   *http.Client
	ApiError interface{}
	*http.Request
}

const (
	GetMethod = "GET"
)

func (c *Client) NewRequest(rawurl string, apierr interface{}) (*Request, error) {
	u, err := c.resolveReferenceString(rawurl)
	if err != nil {
		return nil, err
	}

	httpreq, err := http.NewRequest(GetMethod, u, nil)
	return &Request{c.HttpClient, apierr, httpreq}, err
}
