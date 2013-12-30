// Package httpcache provides facilities for caching HTTP responses and
// hypermedia for REST resources.  The saved responses respect HTTP caching
// policies.  Hypermedia shouldn't change, so is stored for as long as possible.
package httpcache

import (
	"errors"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/hypermedia"
)

var (
	DefaultAdapter  Adapter
	NoResponseError = errors.New("No Response")
)

type Adapter interface {
	Get(string, interface{}) *sawyer.Response
	Set(string, *sawyer.Response, interface{}) error
	Rels(string) hypermedia.Relations
}

// Get retrieves a Response for a REST resource by its URL.  The URL should be
// the full canonical URL for the resource.  The response will be nil if it is
// expired.
func Get(url string, v interface{}) *sawyer.Response {
	return defaultAdapter().Get(url, v)
}

// Set caches a Response for a resource by its URL.
func Set(url string, res *sawyer.Response, v interface{}) error {
	return defaultAdapter().Set(url, res, v)
}

// Rels retrieves the hypmermedia for a REST resource by its URL.  The relations
// cache doesn't expire with the assumption that web services will provide
// redirects as URLs change.
func Rels(url string) hypermedia.Relations {
	return defaultAdapter().Rels(url)
}

func EmptyResponse() *sawyer.Response {
	return sawyer.ResponseError(NoResponseError)
}

func defaultAdapter() Adapter {
	if DefaultAdapter == nil {
		DefaultAdapter = NewMemoryCache()
	}
	return DefaultAdapter
}
