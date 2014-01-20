package httpcache

import (
	"bytes"
	"errors"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/hypermedia"
	"io/ioutil"
	"net/http"
)

type cacheEntry struct {
	Response  *bytes.Reader
	Body      []byte
	Relations hypermedia.Relations
}

// MemoryCache is a sawyer.Cacher that stores the entries in memory.
type MemoryCache struct {
	Cache map[string]*cacheEntry
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{make(map[string]*cacheEntry)}
}

func (c *MemoryCache) Get(req *http.Request) (sawyer.CachedResponse, error) {
	key := RequestKey(req)
	entry, ok := c.Cache[key]
	if !ok {
		return nil, NoResponseError
	}

	cachedResponse, err := Decode(entry.Response)
	if err == nil {
		cachedResponse.Cacher = c
		cachedResponse.SetBodyFunc = func(res *sawyer.Response) {
			res.Body = ioutil.NopCloser(bytes.NewBuffer(entry.Body))
			res.BodyClosed = false
		}
	}

	return cachedResponse, err
}

func (c *MemoryCache) Set(req *http.Request, res *sawyer.Response) error {
	key := RequestKey(req)

	bodyBuffer := &bytes.Buffer{}
	if err := EncodeBody(res, bodyBuffer); err != nil {
		return err
	}

	resBuffer := &bytes.Buffer{}
	if err := Encode(res, resBuffer); err != nil {
		return err
	}

	c.Cache[key] = &cacheEntry{
		bytes.NewReader(resBuffer.Bytes()),
		bodyBuffer.Bytes(),
		nil,
	}

	return nil
}

func (c *MemoryCache) SetRels(req *http.Request, rels hypermedia.Relations) error {
	key := RequestKey(req)
	entry, ok := c.Cache[key]
	if !ok {
		return errors.New("No entry for " + key)
	}

	entry.Relations = rels
	return nil
}

func (c *MemoryCache) Rels(req *http.Request) (hypermedia.Relations, bool) {
	key := RequestKey(req)
	if entry, ok := c.Cache[key]; ok && entry.Relations != nil {
		return entry.Relations, true
	}

	return nil, false
}
