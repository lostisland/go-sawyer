package httpcache

import (
	"bytes"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/hypermedia"
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

func (c *MemoryCache) Get(url string, v interface{}) *sawyer.Response {
	entry := c.Cache[url]
	if entry != nil {
		return DecodeFrom(v, entry.Response, entry.Body)
	}

	return EmptyResponse()
}

func (c *MemoryCache) Set(url string, res *sawyer.Response, v interface{}) error {
	entry := &cacheEntry{&bytes.Buffer{}, &bytes.Buffer{}}
	err := EncodeTo(v, res, entry.Response, entry.Body)
	if err == nil {
		c.Cache[url] = entry
	}
	return err
}

func (c *MemoryCache) Rels(url string) hypermedia.Relations {
	entry := c.Cache[url]
	if entry != nil {
		return Decode(entry.Response).Rels
	}

	return hypermedia.Relations{}
}
