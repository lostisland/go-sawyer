package httpcache

import (
	"bytes"
	"github.com/lostisland/go-sawyer"
	"net/http"
)

type cacheEntry struct {
	Response *bytes.Buffer
	Body     *bytes.Buffer
}

type MemoryCache struct {
	Cache map[string]*cacheEntry
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{make(map[string]*cacheEntry)}
}

func (c *MemoryCache) Get(req *http.Request, v interface{}) *sawyer.Response {
	key := RequestKey(req)
	entry := c.Cache[key]
	if entry != nil {
		return DecodeFrom(v, entry.Response, entry.Body)
	}

	return EmptyResponse()
}

func (c *MemoryCache) Set(req *http.Request, res *sawyer.Response, v interface{}) error {
	key := RequestKey(req)
	entry := &cacheEntry{&bytes.Buffer{}, &bytes.Buffer{}}
	err := EncodeTo(v, res, entry.Response, entry.Body)
	if err == nil {
		c.Cache[key] = entry
	}
	return err
}
