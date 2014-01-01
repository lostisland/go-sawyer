package sawyer

import (
	"errors"
	"github.com/lostisland/go-sawyer/hypermedia"
	"net/http"
)

// See httpcache.Adapter
type Cacher interface {
	Get(*http.Request, interface{}) *Response
	Set(*http.Request, *Response, interface{}) error
	Rels(*http.Request) hypermedia.Relations
}

type NoOpCache struct{}

func (c *NoOpCache) Get(req *http.Request, v interface{}) *Response {
	return noOpResponse
}

func (c *NoOpCache) Set(req *http.Request, res *Response, v interface{}) error {
	return nil
}

func (c *NoOpCache) Rels(req *http.Request) hypermedia.Relations {
	return noOpRels
}

var (
	noOpResponse = ResponseError(errors.New("No Response"))
	noOpRels     = hypermedia.Relations{}
)
