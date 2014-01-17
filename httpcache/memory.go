package httpcache

import (
	"bytes"
	"github.com/lostisland/go-sawyer"
	"io/ioutil"
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

func (c *MemoryCache) Get(req *http.Request) *sawyer.Response {
	key := RequestKey(req)
	entry := c.Cache[key]
	if entry == nil {
		return EmptyResponse()
	}

	res := Decode(entry.Response)
	res.Request = req
	res.Body = ioutil.NopCloser(entry.Body)
	return res
}

func (c *MemoryCache) Set(req *http.Request, res *sawyer.Response) error {
	key := RequestKey(req)
	entry := &cacheEntry{&bytes.Buffer{}, &bytes.Buffer{}}

	err := EncodeBody(res, entry.Body)
	if err != nil {
		return err
	}

	err = Encode(res, entry.Response)
	if err == nil {
		c.Cache[key] = entry
	}
	return err
}
