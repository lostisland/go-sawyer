package httpcache

import (
	"github.com/bmizerany/assert"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/hypermedia"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func CacheResponsesTestFor(cacher sawyer.Cacher, t *testing.T) {
	CacheGet(cacher, t)
	UnusedCache(cacher, t)
	ClearsCache("POST", cacher, t)
	ClearsCache("PUT", cacher, t)
	ClearsCache("PATCH", cacher, t)
	ClearsCache("DELETE", cacher, t)
	GetSetCacheTestFor(cacher, t)
	ETagExpirationTestFor(cacher, t)
}

func CacheGet(cacher sawyer.Cacher, t *testing.T) {
	resp := ""
	srv, cli := server(cacher, func(w http.ResponseWriter, r *http.Request) {
		resp = resp + " "
		w.Header().Set("Content-Length", strconv.Itoa(len(resp)))
		w.WriteHeader(200)
		w.Write([]byte(resp))
	})
	defer srv.Close()

	req, err := cli.NewRequest("/")
	assert.Equal(t, nil, err)

	res := req.Get()
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, int64(1), res.ContentLength)
	by, err := ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, " ", string(by))

	res = req.Get()
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, int64(1), res.ContentLength)
	by, err = ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, " ", string(by))
}

func UnusedCache(cacher sawyer.Cacher, t *testing.T) {
	resp := ""
	srv, cli := server(cacher, func(w http.ResponseWriter, r *http.Request) {
		resp = resp + " "
		w.Header().Set("Content-Length", strconv.Itoa(len(resp)))
		w.WriteHeader(200)
		if r.Method != "HEAD" {
			w.Write([]byte(resp))
		}
	})
	defer srv.Close()

	req, err := cli.NewRequest("/")
	assert.Equal(t, nil, err)

	res := req.Get()
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, int64(1), res.ContentLength)
	by, err := ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, " ", string(by))

	res = req.Options()
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, int64(2), res.ContentLength)
	by, err = ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, "  ", string(by))

	res = req.Head()
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, int64(3), res.ContentLength)

	res = req.Options()
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, int64(4), res.ContentLength)
	by, err = ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, "    ", string(by))

	res = req.Head()
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, int64(5), res.ContentLength)

	res = req.Get()
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, int64(1), res.ContentLength)
	by, err = ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, " ", string(by))
}

func ClearsCache(method string, cacher sawyer.Cacher, t *testing.T) {
	resp := ""
	srv, cli := server(cacher, func(w http.ResponseWriter, r *http.Request) {
		resp = resp + " "
		w.Header().Set("Content-Length", strconv.Itoa(len(resp)))
		w.WriteHeader(200)
		w.Write([]byte(resp))
	})
	defer srv.Close()

	req, err := cli.NewRequest("/")
	assert.Equal(t, nil, err)

	res := req.Get()
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, int64(1), res.ContentLength)
	by, err := ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, " ", string(by))

	res = req.Do(method)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, int64(2), res.ContentLength)
	by, err = ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, "  ", string(by))

	res = req.Get()
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, int64(3), res.ContentLength)
	by, err = ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, "   ", string(by))

	res = req.Get()
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, int64(3), res.ContentLength)
	by, err = ioutil.ReadAll(res.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, "   ", string(by))
}

func GetSetCacheTestFor(cacher sawyer.Cacher, t *testing.T) {
	srv, cli := server(cacher, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"Name":"Resource","Url":"Link"}`))
	})
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

func ETagExpirationTestFor(cacher sawyer.Cacher, t *testing.T) {
	assertExpirationTestFor("ETag", "If-None-Match", `"boom"`, cacher, t)
}

func LastModExpirationTestFor(cacher sawyer.Cacher, t *testing.T) {
	assertExpirationTestFor("Last-Modified", "If-Modified-Since", `"boom"`, cacher, t)
}

func assertExpirationTestFor(reqHeader, resHeader, headerValue string, cacher sawyer.Cacher, t *testing.T) {
	srv, cli := server(cacher, func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		head.Set(reqHeader, headerValue)

		if r.Header.Get(resHeader) != headerValue {
			head.Set("Cache-Control", "max-age=-300")
			w.WriteHeader(200)
			w.Write([]byte(`{"Name":"Resource","Url":"Link"}`))
			return
		}

		w.WriteHeader(304)
		w.Write([]byte(`{"Name":"Changed","Url":"Link"}`))
	})
	defer srv.Close()

	req, err := cli.NewRequest("/")
	assert.Equal(t, nil, err)

	// make initial request
	res := req.Get()
	assert.Equal(t, false, res.IsError())
	assert.Equal(t, 200, res.StatusCode)

	// decode value from initial response
	value := &HttpCacheTestValue{}
	assert.Equal(t, nil, res.Decode(value))
	assert.Equal(t, "Resource", value.Name)

	// cached request is not fresh because max-age=-60
	cached, err := cli.Cacher.Get(req.Request)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, cached.IsFresh())

	// make second request
	res2 := req.Get()
	assert.Equal(t, false, res2.IsError())
	assert.Equal(t, 200, res2.StatusCode)

	// decode value from second response
	value2 := &HttpCacheTestValue{}
	assert.Equal(t, nil, res2.Decode(value2))
	assert.Equal(t, "Resource", value2.Name)

	// after 304, Expiration should be updated
	cached, err = cli.Cacher.Get(req.Request)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, cached.IsFresh())
}

func server(cacher sawyer.Cacher, handler http.HandlerFunc) (*httptest.Server, *sawyer.Client) {
	srv := httptest.NewServer(handler)
	cli, _ := sawyer.NewFromString(srv.URL, nil)
	cli.Cacher = cacher
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
