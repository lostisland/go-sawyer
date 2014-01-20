package httpcache

import (
	"github.com/bmizerany/assert"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/hypermedia"
	"net/http"
	"net/http/httptest"
	"testing"
)

func CacheResponsesTestFor(cacher sawyer.Cacher, t *testing.T) {
	srv, cli := server(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"Name":"Resource","Url":"Link"}`))
	})
	cli.Cacher = cacher
	defer srv.Close()

	req, err := cli.NewRequest("/")
	assert.Equal(t, nil, err)

	// cache is empty
	rels, ok := cli.Cacher.Rels(req.Request)
	assert.Equal(t, false, ok)
	cachedResponse, err := cli.Cacher.Get(req.Request)
	assert.NotEqual(t, nil, err)

	// make first request
	res := req.Get()
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	// response is cached
	cachedResponse, err = cli.Cacher.Get(req.Request)
	assert.Equal(t, nil, err)

	res2 := cachedResponse.Decode(req)
	assert.Equal(t, res.StatusCode, res2.StatusCode)
	assert.Equal(t, res.Header.Get("Content-Type"), res2.Header.Get("Content-Type"))

	// rels are not cached yet
	rels, ok = cli.Cacher.Rels(req.Request)
	assert.Equal(t, false, ok)

	// decode the resource from the original response
	value := &HttpCacheTestValue{}
	rels, ok = value.Rels()
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, res.Decode(value))
	rels, ok = value.Rels()
	assert.Equal(t, true, ok)
	assert.Equal(t, 2, len(rels))
	assert.Equal(t, 2, len(hypermedia.Rels(value)))

	// now rels are cached
	assert.Equal(t, "Resource", value.Name)
	assert.Equal(t, "Link", string(value.Url))
	assert.Equal(t, 2, len(hypermedia.Rels(value)))

	rels, ok = cli.Cacher.Rels(req.Request)
	assert.Equal(t, true, ok)
	assert.Equal(t, 2, len(rels))

	// decode resource from cached response
	value2 := &HttpCacheTestValue{}
	rels, ok = value2.Rels()
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, res2.Decode(value2))
	assert.Equal(t, "Resource", value2.Name)
	assert.Equal(t, "Link", string(value2.Url))
	rels, ok = value2.Rels()
	assert.Equal(t, true, ok)
	assert.Equal(t, 2, len(rels))
	assert.Equal(t, 2, len(hypermedia.Rels(value2)))

	// check that similar request with new accept header is not cached
	req2, err := cli.NewRequest("/")
	assert.Equal(t, nil, err)
	req2.Header.Set("Accept", "application/vnd.sawyer.v2+json")
	_, err = cli.Cacher.Get(req2.Request)
	assert.NotEqual(t, nil, err)

	rels, ok = cli.Cacher.Rels(req2.Request)
	assert.Equal(t, false, ok)
}

func server(handler http.HandlerFunc) (*httptest.Server, *sawyer.Client) {
	srv := httptest.NewServer(handler)
	cli, _ := sawyer.NewFromString(srv.URL, nil)
	cli.Header.Set("Accept", "application.vnd.sawyer+json")
	return srv, cli
}

type HttpCacheTestValue struct {
	Name       string
	Url        hypermedia.Hyperlink
	cachedRels hypermedia.Relations
}

func (r *HttpCacheTestValue) Rels() (hypermedia.Relations, bool) {
	if r.cachedRels == nil {
		return nil, false
	}
	return r.cachedRels, true
}

func (r *HttpCacheTestValue) CacheRels(rels hypermedia.Relations) {
	r.cachedRels = rels
}

func (r *HttpCacheTestValue) HypermediaRels(rels hypermedia.Relations) {
	rels["hypermedia"] = hypermedia.Hyperlink("hypermedia")
}

func (r *HttpCacheTestValue) HyperfieldRels() {}
