package sawyer

import (
	"errors"
	"github.com/lostisland/go-sawyer/hypermedia"
	"net/http"
)

type CachedResource interface {
	Rels() hypermedia.Relations
	CacheRels(hypermedia.Relations)
}

type Cacher interface {
	Get(*http.Request) *Response
	Set(*http.Request, *Response) error
	SetRels(*http.Request, hypermedia.Relations) error
	Rels(*http.Request) hypermedia.Relations
}

type NoOpCache struct{}

func (c *NoOpCache) Get(req *http.Request) *Response {
	return noOpResponse
}

func (c *NoOpCache) Set(req *http.Request, res *Response) error {
	return nil
}

func (c *NoOpCache) SetRels(req *http.Request, rels hypermedia.Relations) error {
	return nil
}

func (c *NoOpCache) Rels(req *http.Request) hypermedia.Relations {
	return make(hypermedia.Relations)
}

var (
	noOpResponse = ResponseError(errors.New("No Response"))
	noOpCacher   Cacher
)

func init() {
	noOpCacher = &NoOpCache{}
}
