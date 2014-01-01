package httpcache

import (
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/hypermedia"
	"net/http"
)

type NoOpCache bool

func (c *NoOpCache) Get(req *http.Request, v interface{}) *sawyer.Response {
	return noOpResponse
}

func (c *NoOpCache) Set(req *http.Request, res *sawyer.Response, v interface{}) error {
	return nil
}

func (c *NoOpCache) Rels(req *http.Request) hypermedia.Relations {
	return noOpRels
}

var (
	noOpResponse = EmptyResponse()
	noOpRels     = hypermedia.Relations{}
)
