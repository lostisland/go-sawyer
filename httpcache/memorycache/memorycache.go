package memorycache

import (
	"bytes"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/httpcache"
)

type cacheEntry struct {
	Response *bytes.Buffer
	Body     *bytes.Buffer
}

type MemoryCache struct {
	Cache map[string]*cacheEntry
}

func New() *MemoryCache {
	return &MemoryCache{make(map[string]*cacheEntry)}
}

func (c *MemoryCache) Get(url string, v interface{}) *sawyer.Response {
	entry := c.Cache[url]
	if entry == nil {
		return nil
	}

	return httpcache.DecodeFrom(v, entry.Response, entry.Body)
}

func (c *MemoryCache) Set(url string, res *sawyer.Response, v interface{}) error {
	entry := &cacheEntry{&bytes.Buffer{}, &bytes.Buffer{}}
	err := httpcache.EncodeTo(v, res, entry.Response, entry.Body)
	if err == nil {
		c.Cache[url] = entry
	}
	return err
}
