// Package httpcache provides facilities for caching HTTP responses and
// hypermedia for REST resources.  The saved responses respect HTTP caching
// policies.  Hypermedia shouldn't change, so is stored for as long as possible.
package httpcache

import (
	"errors"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/hypermedia"
)

// Get retrieves a Response for a REST resource by its URL.  The URL should be
// the full canonical URL for the resource.  The response will be nil if it is
// expired.
func Get(url string) (*sawyer.Response, error) {
	return nil, errors.New("Not Implemented")
}

// Set caches a Response for a resource by its URL.
func Set(url string, res *sawyer.Response) error {
	return errors.New("Not Implemented")

}

// Rels retrieves the hypmermedia for a REST resource by its URL.  The relations
// cache doesn't expire with the assumption that web services will provide
// redirects as URLs change.
func Rels(url string) (hypermedia.Relations, error) {
	return nil, errors.New("Not Implemented")
}
