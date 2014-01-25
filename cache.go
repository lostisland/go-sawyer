package sawyer

import (
	"errors"
	"github.com/lostisland/go-sawyer/hypermedia"
	"net/http"
)

// A Cacher has the ability to get and set caches for HTTP requests and resource
// relations.  See the sawyer/httpcache package.
type Cacher interface {
	// Get gets a CachedResponse for the given request.
	Get(*http.Request) (CachedResponse, error)

	// Set caches the response for the given request.
	Set(*http.Request, *Response) error

	// Reset removes the cached response and body, but leaves the cached relations.
	Reset(*http.Request) error

	// UpdateCache updates the cache for the given request with the expiration from
	// the response.
	UpdateCache(*http.Request, *http.Response) error

	// SetRels caches the given relations for the request.
	SetRels(*http.Request, hypermedia.Relations) error

	// Rels gets the cached relations for the given request.
	Rels(*http.Request) (hypermedia.Relations, bool)
}

// CachedResponse is an interface for the httpcache CachedResponseDecoder.
type CachedResponse interface {
	Decode(*Request) *Response
	SetupRequest(*http.Request)
	IsFresh() bool
	IsExpired() bool
}

type noOpCache struct{}

func (c *noOpCache) Get(req *http.Request) (CachedResponse, error) {
	return nil, noOpError
}

func (c *noOpCache) Set(req *http.Request, res *Response) error {
	return nil
}

func (c *noOpCache) UpdateCache(req *http.Request, res *http.Response) error {
	return nil
}

func (c *noOpCache) Reset(req *http.Request) error {
	return nil
}

func (c *noOpCache) SetRels(req *http.Request, rels hypermedia.Relations) error {
	return nil
}

func (c *noOpCache) Rels(req *http.Request) (hypermedia.Relations, bool) {
	return nil, false
}

var (
	noOpError  = errors.New("No Response")
	noOpCacher Cacher
)

func init() {
	noOpCacher = &noOpCache{}
}
