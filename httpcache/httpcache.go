// Package httpcache provides facilities for caching HTTP responses and
// hypermedia for REST resources.  The saved responses respect HTTP caching
// policies.  Hypermedia shouldn't change, so is stored for as long as possible.
package httpcache

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/lostisland/go-sawyer"
	"net/http"
)

// RequestKey builds a unique string key for a net/http Request.
func RequestKey(r *http.Request) string {
	return r.Header.Get(keyHeader) + keySep + r.URL.String()
}

func RequestSha(r *http.Request) string {
	key := RequestKey(r)
	sum := sha256.New().Sum([]byte(key))
	return hex.EncodeToString(sum)
}

func ResponseError(err error) *sawyer.Response {
	return sawyer.ResponseError(err)
}

var NoResponseError = errors.New("No Response")

const (
	keySep    = ":"
	keyHeader = "Accept"
)
