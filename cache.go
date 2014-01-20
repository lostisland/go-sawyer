package sawyer

import (
	"errors"
	"github.com/lostisland/go-sawyer/hypermedia"
)

// See httpcache.Adapter
type Cacher interface {
	Get(*Request) *Response
	Set(*Request, *Response) error
	SetRels(*Request, hypermedia.Relations) error
	Rels(*Request) hypermedia.Relations
}

type NoOpCache struct{}

func (c *NoOpCache) Get(req *Request) *Response {
	return noOpResponse
}

func (c *NoOpCache) Set(req *Request, res *Response) error {
	return nil
}

func (c *NoOpCache) SetRels(req *Request, rels hypermedia.Relations) error {
	return nil
}

func (c *NoOpCache) Rels(req *Request) hypermedia.Relations {
	return make(hypermedia.Relations)
}

var noOpResponse = ResponseError(errors.New("No Response"))
