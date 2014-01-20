// Package httpcache provides facilities for caching HTTP responses and
// hypermedia for REST resources.  The saved responses respect HTTP caching
// policies.  Hypermedia shouldn't change, so is stored for as long as possible.
package httpcache

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/hypermedia"
)

var (
	DefaultAdapter  Adapter
	NoResponseError = errors.New("No Response")
)

type Adapter interface {
	// Get retrieves a Response for a REST resource by its URL.  The URL should be
	// the full canonical URL for the resource.  The response will be nil if it is
	// expired.
	Get(*sawyer.Request) *sawyer.Response

	// Set caches a Response for a resource by its URL.
	Set(*sawyer.Request, *sawyer.Response) error

	SetRels(*sawyer.Request, hypermedia.Relations) error

	Rels(*sawyer.Request) hypermedia.Relations
}

func ResponseError(err error) *sawyer.Response {
	return sawyer.ResponseError(err)
}

func EmptyResponse() *sawyer.Response {
	return ResponseError(NoResponseError)
}

func RequestKey(r *sawyer.Request) string {
	return r.Header.Get(keyHeader) + keySep + r.URL.String()
}

func RequestSha(r *sawyer.Request) string {
	key := RequestKey(r)
	sum := sha256.New().Sum([]byte(key))
	return hex.EncodeToString(sum)
}

const (
	keySep    = ":"
	keyHeader = "Accept"
)
