package httpcache

import (
	"github.com/bmizerany/assert"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/hypermedia"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCacheResponses(t *testing.T) {
	srv, cli := server(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"Name":"Resource","Url":"Link"}`))
	})
	defer srv.Close()

	req, err := cli.NewRequest("/")
	assert.Equal(t, nil, err)

	// cache is empty
	assert.Equal(t, true, cli.Cacher.Get(req).IsError())
	assert.Equal(t, 0, len(cli.Cacher.Rels(req)))

	res := req.Get()
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	// rels are not cached yet
	assert.Equal(t, 0, len(cli.Cacher.Rels(req)))

	value := &HttpCacheTestValue{}
	assert.Equal(t, nil, res.Decode(value))

	assert.Equal(t, "Resource", value.Name)
	assert.Equal(t, "Link", string(value.Url))
	assert.Equal(t, 2, len(hypermedia.Rels(value)))
	//assert.Equal(t, 2, len(cli.Cacher.Rels(req)))

	// response is cached
	res2 := cli.Cacher.Get(req)
	assert.Equal(t, false, res2.IsError())
	res2 = cli.Cacher.Get(req)
	assert.Equal(t, false, res2.IsError())
	assert.Equal(t, res.StatusCode, res2.StatusCode)
	assert.Equal(t, res.Header.Get("Content-Type"), res2.Header.Get("Content-Type"))

	value2 := &HttpCacheTestValue{}
	assert.Equal(t, nil, res2.Decode(value2))
	assert.Equal(t, "Resource", value2.Name)
	assert.Equal(t, "Link", string(value2.Url))
	assert.Equal(t, 2, len(hypermedia.Rels(value2)))

	req2, err := cli.NewRequest("/")
	assert.Equal(t, nil, err)
	req2.Header.Set("Accept", "application/vnd.sawyer.v2+json")
	assert.Equal(t, true, cli.Cacher.Get(req2).IsError())
	assert.Equal(t, 0, len(cli.Cacher.Rels(req2)))
}

func server(handler http.HandlerFunc) (*httptest.Server, *sawyer.Client) {
	srv := httptest.NewServer(handler)
	cli, _ := sawyer.NewFromString(srv.URL, nil)
	cli.Header.Set("Accept", "application.vnd.sawyer+json")
	cli.Cacher = NewMemoryCache()
	return srv, cli
}

type HttpCacheTestValue struct {
	Name string
	Url  hypermedia.Hyperlink
}

func (r *HttpCacheTestValue) HypermediaRels(rels hypermedia.Relations) {
	rels["hypermedia"] = hypermedia.Hyperlink("hypermedia")
}

func (r *HttpCacheTestValue) HyperfieldRels() {}
