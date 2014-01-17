package sawyer

import (
	"errors"
	"net/http"
)

// See httpcache.Adapter
type Cacher interface {
	Get(*http.Request) *Response
	Set(*http.Request, *Response) error
}

type NoOpCache struct{}

func (c *NoOpCache) Get(req *http.Request) *Response {
	return noOpResponse
}

func (c *NoOpCache) Set(req *http.Request, res *Response) error {
	return nil
}

var noOpResponse = ResponseError(errors.New("No Response"))
