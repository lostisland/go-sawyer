package httpcache

import (
	"bytes"
	"errors"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/hypermedia"
	"io/ioutil"
	"net/http"
)

// MemoryCache is a sawyer.Cacher that stores the entries in memory.  This is
// only intended for testing, and should not be used in production.
type MemoryCache struct {
	Cache map[string]*cacheEntry
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{make(map[string]*cacheEntry)}
}

func (c *MemoryCache) Get(req *http.Request) (sawyer.CachedResponse, error) {
	if _, entry, ok := c.getEntry(req); ok {
		return entry.Decode(c)
	}

	return nil, NoResponseError
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

func (c *MemoryCache) Reset(req *http.Request) error {
	if key, entry, ok := c.getEntry(req); ok {
		entry.Response = nil
		c.Cache[key] = entry
	}

	return nil
}

func (c *MemoryCache) UpdateCache(req *http.Request, res *http.Response) error {
	key, entry, ok := c.getEntry(req)
	if !ok {
		return NoResponseError
	}

	cached, err := entry.Decode(c)
	if err != nil {
		return err
	}

	cached.Expires = expiration(res)

	buf := &bytes.Buffer{}
	EncodeResponse(cached.CachedResponse, buf)
	entry.Response = bytes.NewReader(buf.Bytes())
	c.Cache[key] = entry
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

type cacheEntry struct {
	Response  *bytes.Reader
	Body      []byte
	Relations hypermedia.Relations
}

func (e *cacheEntry) Decode(cacher sawyer.Cacher) (*CachedResponseDecoder, error) {
	cachedResponse, err := Decode(e.Response)
	e.Response.Seek(0, 0)

	if err == nil {
		cachedResponse.Cacher = cacher
		cachedResponse.SetBodyFunc = func(res *sawyer.Response) {
			res.Body = ioutil.NopCloser(bytes.NewBuffer(e.Body))
			res.BodyClosed = false
		}
	}

	return cachedResponse, err
}

func (c *MemoryCache) getEntry(req *http.Request) (string, *cacheEntry, bool) {
	key := RequestKey(req)
	entry, ok := c.Cache[key]
	if ok && entry.Response == nil {
		ok = false
	}
	return key, entry, ok
}
