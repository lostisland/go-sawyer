package sawyer

import (
	"errors"
	"github.com/lostisland/go-sawyer/hypermedia"
	"net/http"
)

// A Cacher has the ability to get and set caches for HTTP requests and resource
// relations.  See the sawyer/httpcache package.
type Cacher interface {
	Get(*http.Request) *Response
	Set(*http.Request, *Response) error
	SetRels(*http.Request, hypermedia.Relations) error
	Rels(*http.Request) (hypermedia.Relations, bool)
}

type noOpCache struct{}

func (c *noOpCache) Get(req *http.Request) *Response {
	return noOpResponse
}

func (c *noOpCache) Set(req *http.Request, res *Response) error {
	return nil
}

func (c *noOpCache) SetRels(req *http.Request, rels hypermedia.Relations) error {
	return nil
}

func (c *noOpCache) Rels(req *http.Request) (hypermedia.Relations, bool) {
	return nil, false
}

var (
	noOpResponse = ResponseError(errors.New("No Response"))
	noOpCacher   Cacher
)

func init() {
	noOpCacher = &noOpCache{}
}
