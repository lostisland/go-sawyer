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

	res := httpcache.Decode(entry.Response)

	if v != nil && res.ContentLength > 0 {
		err := res.MediaType.Decode(v, entry.Body)
		if err != nil {
			return sawyer.ResponseError(err)
		}
	}

	return res
}
