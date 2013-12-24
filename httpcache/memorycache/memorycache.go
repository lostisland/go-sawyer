package memorycache

import (
	"bytes"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/httpcache"
	"io"
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

func (c *MemoryCache) Set(url string, res *sawyer.Response, v interface{}) error {
	entry := &cacheEntry{}

	if v != nil && res.ContentLength > 0 {
		entry.Body = &bytes.Buffer{}
		reader := io.TeeReader(res.Body, entry.Body)
		dec, err := res.MediaType.Decoder(reader)
		if err != nil {
			return err
		}

		err = dec.Decode(v)
		if err != nil {
			return err
		}
	}

	entry.Response = &bytes.Buffer{}
	err := httpcache.Encode(res, entry.Response)
	if err != nil {
		return err
	}

	c.Cache[url] = entry
	return nil
}
